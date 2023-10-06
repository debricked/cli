package files

import (
	"github.com/debricked/cli/internal/cmd/files/find"
	"github.com/debricked/cli/internal/cmd/files/fingerprint"
	"github.com/debricked/cli/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewFilesCmd(finder file.IFinder, fingerprinter file.IFingerprint) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(find.NewFindCmd(finder))
	cmd.AddCommand(fingerprint.NewFingerprintCmd(fingerprinter))

	return cmd
}
