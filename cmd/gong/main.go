package main

import (
	"fmt"
	"github.com/kensodev/gong"
	"github.com/segmentio/go-prompt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login to your Jira Instance",
			Action: func(c *cli.Context) error {
				username := prompt.String("What is your Jira username/email?")
				password := prompt.PasswordMasked("What is your password?")
				domain := prompt.String("What is the jira instance URL?")

				loginDetails := gong.NewLoginDetails(username, password, domain)
				err := loginDetails.Verify()

				if err != nil {
					fmt.Println("Unable to login, please check your credentials")
					return err
				}

				err = loginDetails.Save()

				if err != nil {
					fmt.Println("Unable to save file to disk")
					return err
				}

				fmt.Println("Successfully authenticated!, saved login details to disk")

				return nil
			},
		},
	}
	app.Run(os.Args)

}
