package main

import (
	"fmt"
	"github.com/kensodev/gong"
	"github.com/segmentio/go-prompt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
)

func main() {
	var branchType string

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
				_, err := loginDetails.GetClient()

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
		{
			Name:  "start",
			Usage: "Start working on a ticket. Creates a branch on your local repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "type",
					Value:       "feature",
					Usage:       "Type of branch to create eg: feature/{ticket-id}-ticket-title",
					Destination: &branchType,
				},
			},
			Action: func(c *cli.Context) error {
				issueId := c.Args()[0]
				jiraClient, err := gong.GetAuthenticatedClient()
				if err != nil {
					fmt.Println(err)
					return err
				}

				branchName := gong.GetBranchName(jiraClient, issueId, branchType)

				cmd := "git"
				args := []string{"checkout", "-b", branchName}

				out, err := exec.Command(cmd, args...).Output()

				if err != nil {
					fmt.Println(err)
					return err
				}

				fmt.Println(string(out))

				return nil
			},
		},
	}
	app.Run(os.Args)
}
