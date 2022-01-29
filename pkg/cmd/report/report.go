package report

import (
	"debricked/pkg/client"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var debClient *client.DebClient

func Report() error {
	fmt.Println("Reporting!")
	return nil
}

func NewReportCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	return &cobra.Command{
		Use:   "report",
		Short: "Generate license report for upload ID",
		Long:  "Generate license report for upload ID",
		Run: func(cmd *cobra.Command, args []string) {
			err := Report()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
