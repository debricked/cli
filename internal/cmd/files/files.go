package files

import (
	"github.com/debricked/cli/pkg/cmd/files/find"
	"github.com/debricked/cli/pkg/cmd/files/resolve"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/resolution"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewFilesCmd(finder file.IFinder, resolver resolution.IResolver) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(find.NewFindCmd(finder))
	cmd.AddCommand(resolve.NewResolveCmd(resolver))

	return cmd
}
