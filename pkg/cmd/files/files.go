package files

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/files/find"
	"debricked/pkg/file"
	"github.com/spf13/cobra"
)

func NewFilesCmd(debClient *client.IDebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Analyze files",
		Long:  "Analyze files",
	}

	f, _ := file.NewFinder(*debClient)
	cmd.AddCommand(find.NewFindCmd(f))

	return cmd
}
