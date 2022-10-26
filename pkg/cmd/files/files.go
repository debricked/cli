package files

import (
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/cmd/files/find"
	"github.com/debricked/cli/pkg/file"
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

	return cmd
}
