package main

import (
	"github.com/fatih/color"
	"github.com/kensodev/gong"
	"github.com/urfave/cli"
	"os"
	"os/exec"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"

	var branchType string

	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login to your Jira Instance",
			Action: func(c *cli.Context) error {
				clientName := c.Args()[0]
				client, err := gong.NewClient(clientName)

				if err != nil {
					color.Red(err.Error())
					return nil
				}

				_, err = gong.Login(client)

				if err != nil {
					color.Red(err.Error())
					return err
				}

				color.Green("Logged in!")

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
				if c.NArg() == 0 {
					color.Red("You have to pass in a ticket id as an argument")
					return nil
				}

				issueId := c.Args()[0]

				client, err := gong.NewAuthenticatedClient()

				if err != nil {
					color.Red("Problem with starting the issue")
				}

				branchName, err := gong.Start(client, branchType, issueId)

				if err != nil {
					color.Red("Problem with starting the issue")
				}

				cmd := "git"
				args := []string{"checkout", "-b", branchName}

				out, err := exec.Command(cmd, args...).Output()

				if err != nil {
					color.Red(err.Error())
					return err
				}

				color.Green(string(out))

				return nil
			},
		},
	}
	app.Run(os.Args)
}
