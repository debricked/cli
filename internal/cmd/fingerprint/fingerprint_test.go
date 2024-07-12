package fingerprint

import (
	"os"
	"testing"

	"github.com/debricked/cli/internal/fingerprint"
	"github.com/debricked/cli/internal/fingerprint/testdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewFingerprintCmd(t *testing.T) {
	var f fingerprint.IFingerprint
	cmd := NewFingerprintCmd(f)

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
		ExclusionFlag,
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

func TestRunE(t *testing.T) {
	defer func() {
		os.Remove(fingerprint.OutputFileNameFingerprints) // TODO: make sure it will remove all generated fingerprint files (e.g. with date suffix)
	}()
	fingerprintMock := testdata.NewFingerprintMock()
	runE := RunE(fingerprintMock)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
	// TODO: Run command again, first with regenerate=true (default) to ensure no extra file is generated (check only one exists)
	// 		 Run command a third time, without regenerate,
	//		 asserting that a date-stamped debricked fingerprint file is created when regenerate=false

}
