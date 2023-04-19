// Package command provides system command execution.
package command

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	v2 "github.com/sensu/core/v2"
	bytesutil "github.com/sensu/sensu-go/util/bytes"
	"github.com/sirupsen/logrus"
)

const undocumentedTestCheckCommand = "!sensu_test_check!"

const cannedResponseText = `
                         .'loo:,
                        ,KNMMWNWX
                  ..    ,000OkxkW'
                 ,o,.   .O0KOOOk0:
                 dkl.    :OO0kkkk;
                 :ko     .lOOOOkx
              'OXX0:       'xWMdkd;o,.
              cMMMM; .,lkXN;oWNOk,;MMMMWXx:
              oMMMM:KNWMMMMl'.cl .NMMMMMMMMX
              NMMMWkMMMMMMMMWxxKONMMMMMMMMMM
             oMMMMMMMMMMMMMMMNW0NMMMMMMMMMMM.
             KMMMMMMMMMMMMMMWMWWMMMMMMMMMMMMN
             oKXXKKKKXMMMMMMMMMWMMMMMMMMMMMMM.
                     'MMMMMMMMMMMMMMMMMMMMMMMk
                     .MMMMMMMMWMMMMMMMMMMMMMMM
                      WMMMMMMMMMMWNMMMMMMMMMMN
                      WMMMMMMMMMWX0kO0WMMMMMMO
                     .MMMMMMMMMMMNX0kkWMMMMWO'
                     ;MMMMMMMMMMMMWXNNMMMMW.
`

var CannedResponse = &ExecutionResponse{
	Output: cannedResponseText,
}

// Executor ...
type Executor interface {
	Execute(context.Context, ExecutionRequest) (*ExecutionResponse, error)
}

const (
	// TimeoutOutput specifies the command execution output in the
	// event of an execution timeout.
	TimeoutOutput string = "Execution timed out\n"

	// OKExitStatus specifies the command execution exit status
	// that indicates a success, A-OK.
	OKExitStatus int = 0

	// TimeoutExitStatus specifies the command execution exit
	// status in the event of an execution timeout.
	TimeoutExitStatus int = 2

	// FallbackExitStatus specifies the command execution exit
	// status used when golang is unable to determine the exit
	// status.
	FallbackExitStatus int = 3
)

// ExecutionRequest provides information about a system command execution,
// somewhat of an abstraction intended to be used for Sensu check,
// mutator, and handler execution.
type ExecutionRequest struct {
	// Command is the command to be executed.
	Command string

	// Env ...
	Env []string

	// Input to provide the command via STDIN.
	Input string

	// Execution timeout in seconds, will be set to a default if
	// not specified.
	Timeout int

	// Name is the name of the resource that is invoking the execution.
	Name string

	// InProgress is a map of checks that are still in execution, this is
	// necessary for a check or hook to escape zombie processes.
	InProgress map[string]*v2.CheckConfig

	// InProgressMu is the mutex for the InProgress map.
	InProgressMu *sync.Mutex
}

// ExecutionResponse provides the response information of an ExecutionRequest.
type ExecutionResponse struct {
	// Combined command execution STDOUT/ERR.
	Output string

	// Command execution exit status.
	Status int

	// Duration provides command execution time in seconds.
	Duration float64
}

// NewExecutor ...
func NewExecutor() Executor {
	return &ExecutionRequest{}
}

// Execute executes a system command (fork/exec) with a
// timeout, optionally writing to STDIN, capturing its combined output
// (STDOUT/ERR) and exit status.
func (e *ExecutionRequest) Execute(ctx context.Context, execution ExecutionRequest) (*ExecutionResponse, error) {
	if execution.Command == undocumentedTestCheckCommand {
		return CannedResponse, nil
	}
	resp := &ExecutionResponse{}
	logger := logrus.WithFields(logrus.Fields{"component": "command"})
	// Using a platform specific shell to "cheat", as the shell
	// will handle certain failures for us, where golang exec is
	// known to have troubles, e.g. command not found. We still
	// use a fallback exit status in the unlikely event that the
	// exit status cannot be determined.
	var cmd *exec.Cmd

	// Use context.WithCancel for command execution timeout.
	// context.WithTimeout will not kill child/grandchild processes
	// (see issues tagged in https://github.com/sensu/sensu-go/issues/781).
	// Rather, we will use a timer, CancelFunc and proc functions
	// to perform full cleanup.
	ctx, timeout := context.WithCancel(ctx)
	defer timeout()

	// Taken from Sensu-Spawn (Sensu 1.x.x).
	cmd = Command(ctx, execution.Command)

	// Set the ENV for the command if it is set
	if len(execution.Env) > 0 {
		cmd.Env = execution.Env
	}

	// Share an output buffer between STDOUT/ERR, following the
	// Nagios plugin spec.
	var output bytesutil.SyncBuffer

	cmd.Stdout = &output
	cmd.Stderr = &output

	// If Input is specified, write to STDIN.
	if execution.Input != "" {
		cmd.Stdin = strings.NewReader(execution.Input)
	}

	started := time.Now()
	defer func() {
		resp.Duration = time.Since(started).Seconds()
	}()

	timer := time.NewTimer(math.MaxInt64)
	defer timer.Stop()
	if execution.Timeout != 0 {
		SetProcessGroup(cmd)
		timer.Stop()
		timer = time.NewTimer(time.Duration(execution.Timeout) * time.Second)
	}
	if err := cmd.Start(); err != nil {
		// Something unexpected happened when attempting to
		// fork/exec, return immediately.
		return resp, err
	}

	waitCh := make(chan struct{})
	var err error
	go func() {
		err = cmd.Wait()
		close(waitCh)
	}()

	// Wait for the process to complete or the timer to trigger, whichever comes first.
	var killErr error
	select {
	case <-waitCh:
		resp.Output = output.String()
		if err != nil {
			// The command most likely return a non-zero exit status.
			if exitError, ok := err.(*exec.ExitError); ok {
				// Best effort to determine the exit status, this
				// should work on Linux, OSX, and Windows.
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					resp.Status = status.ExitStatus()
				} else {
					resp.Status = FallbackExitStatus
				}
			} else {
				resp.Status = FallbackExitStatus
			}
		} else {
			// Everything is A-OK.
			resp.Status = OKExitStatus
		}

	case <-timer.C:
		var killErrOutput string
		if killErr = KillProcess(cmd); killErr != nil {
			logger.WithError(killErr).Errorf("Execution timed out - Unable to TERM/KILL the process: #%d", cmd.Process.Pid)
			killErrOutput = fmt.Sprintf("Unable to TERM/KILL the process: #%d\n", cmd.Process.Pid)
			escapeZombie(&execution)
		}
		timeout()
		resp.Output = fmt.Sprintf("%s%s%s", TimeoutOutput, killErrOutput, output.String())
		resp.Status = TimeoutExitStatus
	}

	return resp, nil
}

func escapeZombie(ex *ExecutionRequest) {
	logger := logrus.WithFields(logrus.Fields{"component": "command"})
	if ex.InProgress != nil && ex.InProgressMu != nil && ex.Name != "" {
		logger.WithField("check", ex.Name).Warn("check or hook execution created zombie process - escaping in order for the check to execute again")
		ex.InProgressMu.Lock()
		delete(ex.InProgress, ex.Name)
		ex.InProgressMu.Unlock()
	} else {
		logger.Error("unable to escape zombie process created from command execution")
	}
}

// FixtureExecutionResponse returns an Execution for use in testing
func FixtureExecutionResponse(status int, output string) *ExecutionResponse {
	return &ExecutionResponse{
		Output:   output,
		Status:   status,
		Duration: 1,
	}
}
