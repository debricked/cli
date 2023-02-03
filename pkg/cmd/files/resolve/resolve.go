package resolve

import (
	"github.com/debricked/cli/pkg/resolution"
	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/strategy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve manifest files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(),
	}

	return cmd
}

func RunE() func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {

		resolver := resolution.NewResolver(
			file.NewBatchFactory(),
			strategy.NewStrategyFactory(),
			resolution.NewScheduler(),
		)

		_, err := resolver.Resolve(args)

		return err
	}
}
