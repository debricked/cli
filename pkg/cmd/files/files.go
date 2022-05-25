package files

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/files/find"
	"github.com/spf13/cobra"
)

func NewFilesCmd(debClient *client.DebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
	}

	cmd.AddCommand(find.NewFindCmd(debClient))

	return cmd
}
