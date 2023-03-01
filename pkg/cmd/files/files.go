package files

import (
	"github.com/debricked/cli/pkg/cmd/files/find"
	"github.com/debricked/cli/pkg/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewFilesCmd(finder file.IFinder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(find.NewFindCmd(finder))

	return cmd
}
