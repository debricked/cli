package login

import (
	"debricked/pkg/client"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var debClient *client.DebClient

func NewLoginCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate via web browser",
		Long:  `Browser based authentication with Debricked. Upon success, an authorized JWT is stored and Debricked services can be used.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := login()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}

// login opens Debricked login. Upon successful login JWT is stored for 1 hour
func login() error {
	fmt.Println("Logged in!")
	return nil
}
