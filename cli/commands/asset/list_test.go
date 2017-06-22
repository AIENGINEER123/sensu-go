package asset

import (
	"errors"
	"testing"

	"github.com/sensu/sensu-go/cli"
	client "github.com/sensu/sensu-go/cli/client/testing"
	test "github.com/sensu/sensu-go/cli/commands/testing"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	cmd := ListCommand(cli)

	assert.NotNil(cmd, "cmd should be returned")
	assert.NotNil(cmd.RunE, "cmd should be able to be executed")
	assert.Regexp("list", cmd.Use)
	assert.Regexp("assets", cmd.Short)
}

func TestListCommandRunEClosureWithTable(t *testing.T) {
	assert := assert.New(t)
	cli := newCLI()

	config := cli.Config.(*client.MockConfig)
	config.On("Format").Return("none")

	client := cli.Client.(*client.MockClient)
	client.On("ListAssets").Return([]types.Asset{
		*types.FixtureAsset("one"),
		*types.FixtureAsset("two"),
	}, nil)

	cmd := ListCommand(cli)
	out, err := test.RunCmd(cmd, []string{})

	assert.NotEmpty(out)
	assert.Nil(err)
}

func TestListCommandRunEClosureWithJSON(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	client := cli.Client.(*client.MockClient)
	client.On("ListAssets").Return([]types.Asset{
		*types.FixtureAsset("one"),
		*types.FixtureAsset("two"),
	}, nil)

	cmd := ListCommand(cli)
	out, err := test.RunCmd(cmd, []string{})

	assert.NotEmpty(out)
	assert.Nil(err)
}

func TestListCommandRunEClosureWithErr(t *testing.T) {
	assert := assert.New(t)

	cli := newCLI()
	client := cli.Client.(*client.MockClient)
	client.On("ListAssets").Return([]types.Asset{}, errors.New("fire"))

	cmd := ListCommand(cli)
	out, err := test.RunCmd(cmd, []string{})

	assert.Empty(out)
	assert.NotNil(err)
	assert.Equal("fire", err.Error())
}

func newCLI() *cli.SensuCli {
	cli := test.NewMockCLI()
	config := cli.Config.(*client.MockConfig)
	config.On("Format").Return("json")

	return cli
}
