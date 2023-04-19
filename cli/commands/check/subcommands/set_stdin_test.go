package subcommands

import (
	"fmt"
	"testing"

	v2 "github.com/sensu/core/v2"
	client "github.com/sensu/sensu-go/cli/client/testing"
	test "github.com/sensu/sensu-go/cli/commands/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetSTDINCommand(t *testing.T) {
	testCases := []struct {
		testName       string
		args           []string
		fetchResponse  error
		updateResponse error
		expectedOutput string
		expectError    bool
	}{
		{"no args", []string{}, nil, nil, "Usage", true},
		{"fetch error", []string{"checky", "foo"}, fmt.Errorf("error"), nil, "", true},
		{"update error", []string{"checky", "bar"}, nil, fmt.Errorf("error"), "", true},
		{"invalid input", []string{"checky", "yes"}, nil, nil, "", true},
		{"valid input", []string{"checky", "true"}, nil, nil, "Updated", false},
	}

	for _, tc := range testCases {
		var name string
		if len(tc.args) > 0 {
			name = tc.args[0]
		}

		t.Run(tc.testName, func(t *testing.T) {
			check := v2.FixtureCheckConfig("checky")
			cli := test.NewMockCLI()

			client := cli.Client.(*client.MockClient)
			client.On(
				"FetchCheck",
				name,
			).Return(check, tc.fetchResponse)

			client.On(
				"UpdateCheck",
				mock.Anything,
			).Return(tc.updateResponse)

			cmd := SetSTDINCommand(cli)
			out, err := test.RunCmd(cmd, tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Regexp(t, tc.expectedOutput, out)
		})
	}
}
