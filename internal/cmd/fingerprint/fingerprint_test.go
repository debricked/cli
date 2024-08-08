package fingerprint

import (
	"io"
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
		os.Remove(fingerprint.OutputFileNameFingerprints)
	}()
	fingerprintMock := testdata.NewFingerprintMock()
	runE := RunE(fingerprintMock)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunEFileExistsError(t *testing.T) {
	defer func() {
		os.Remove(fingerprint.OutputFileNameFingerprints)
	}()
	fingerprintMock := testdata.NewFingerprintMockFileExistsError()
	runE := RunE(fingerprintMock)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runE(nil, []string{"."})

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.NoError(t, err)
	assert.Contains(t, string(output), "change flag '--regenerate' to 'true'")
}
