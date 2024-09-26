package sbom

import (
	"fmt"

	"github.com/debricked/cli/internal/report"
	"github.com/debricked/cli/internal/report/sbom"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitId string
var repositoryId string

const CommitFlag = "commit"
const RepositorylFlag = "repository"
const TokenFlag = "token"

func NewSBOMCmd(reporter report.IReporter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sbom",
		Short: "Generate SBOM report",
		Long: `Generate SBOM report for chosen commit and repository. 
This is an enterprise feature. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(reporter),
	}

	cmd.Flags().StringVarP(&commitId, CommitFlag, "c", "", "The commit that you want an SBOM report for")
	_ = cmd.MarkFlagRequired(CommitFlag)
	viper.MustBindEnv(CommitFlag)
	cmd.Flags().StringVarP(&repositoryId, RepositorylFlag, "r", "", "The repository that you want an SBOM report for")
	_ = cmd.MarkFlagRequired(RepositorylFlag)
	viper.MustBindEnv(RepositorylFlag)

	return cmd
}

func RunE(r report.IReporter) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		orderArgs := sbom.OrderArgs{
			RepositoryID: viper.GetString(RepositorylFlag),
			CommitID:     viper.GetString(CommitFlag),
		}

		if err := r.Order(orderArgs); err != nil {
			return fmt.Errorf("%s %s", color.RedString("⨯"), err.Error())
		}

		return nil
	}
}
