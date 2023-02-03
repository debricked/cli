package files

import (
	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/cmd/files/find"
	"github.com/debricked/cli/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewFilesCmd(debClient *client.IDebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	f, _ := file.NewFinder(*debClient)
	cmd.AddCommand(find.NewFindCmd(f))

	cmd.AddCommand(resolve.NewResolveCmd())

	return cmd
}
