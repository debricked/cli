package resolve

import (
	"github.com/debricked/cli/pkg/resolution"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewResolveCmd(resolver resolution.IResolver) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve manifest files",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(resolver),
	}

	return cmd
}

func RunE(resolver resolution.IResolver) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		_, err := resolver.Resolve(args)

		return err
	}
}
