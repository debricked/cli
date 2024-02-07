package experience

import (
	"github.com/debricked/cli/internal/experience"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ExclusionFlag = "exclusion-experience"
)

func NewExperienceCmd(experienceCalculator experience.IExperience) *cobra.Command {

	short := "Experience calculator uses git blame and call graphs to calculate who has written code with what open source. [beta feature]"
	cmd := &cobra.Command{
		Use:    "xp [path]",
		Short:  short,
		Hidden: false,
		Long:   short, //TODO: Add long description
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(experienceCalculator),
	}

	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func RunE(e experience.IExperience) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		output, err := e.CalculateExperience(path, viper.GetStringSlice(ExclusionFlag))

		if err != nil {
			return err
		}

		err = output.ToFile(experience.OutputFileNameExperience)
		if err != nil {
			return err
		}

		return nil
	}
}
