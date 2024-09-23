package hello

import (
	"testing"

	"github.com/debricked/cli/internal/hello"
	"github.com/debricked/cli/internal/hello/testdata"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewHelloCmd(t *testing.T) {
	cmd := NewHelloCmd(hello.NewDebrickedGreeter())
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equal(t, shorthand, flag.Shorthand)
	}

	var flagKeys = []string{
		NameFlag,
	}
	viperKeys := viper.AllKeys()
	for _, flagKey := range flagKeys {
		match := false
		for _, key := range viperKeys {
			if key == flagKey {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that flag was present: "+flagKey)
	}
}

func TestPreRun(t *testing.T) {
	mockAuthenticator := testdata.MockGreeter{}
	cmd := NewHelloCmd(mockAuthenticator)
	cmd.PreRun(cmd, nil)
}

func TestRunE(t *testing.T) {
	a := testdata.MockGreeter{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}
