package hello

import (
	"fmt"

	"github.com/debricked/cli/internal/hello"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var name string

const NameFlag = "name"

func NewHelloCmd(g hello.IGreeter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Long:  `Says hello to you.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(g),
	}
	cmd.Flags().StringVarP(&name, NameFlag, "n", "anon", `The name to greet`)
	viper.MustBindEnv(NameFlag)

	return cmd
}

func RunE(g hello.IGreeter) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		name := viper.GetString(NameFlag)
		greeting := g.Greeting(name)
		fmt.Println(greeting)

		return nil
	}
}
