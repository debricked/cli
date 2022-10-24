package files

import (
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/cmd/files/find"
	"github.com/debricked/cli/pkg/file"
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
