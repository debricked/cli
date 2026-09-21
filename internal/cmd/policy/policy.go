package policy

import (
	"github.com/debricked/cli/internal/cmd/policy/validate"
	"github.com/debricked/cli/internal/policy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewPolicyCmd(validator policy.IValidator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage Debricked policy files.",
		Long:  `Work with Debricked policy files, such as validating them before a scan.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(validate.NewValidateCmd(validator))

	return cmd
}
