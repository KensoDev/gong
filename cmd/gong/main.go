package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/KensoDev/gong"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.4.0"

	var branchType string

	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login to your project managment tool instance",
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
		{
			Name:  "browse",
			Usage: "Browse to the jira URL of the branch you are working on",
			Action: func(c *cli.Context) error {
				cmd := "git"
				args := []string{"rev-parse", "--abbrev-ref", "HEAD"}

				out, err := exec.Command(cmd, args...).Output()

				if err != nil {
					return err
				}

				branchName := string(out)

				client, err := gong.NewAuthenticatedClient()

				if err != nil {
					color.Red(err.Error())
					return err
				}

				url, err := gong.Browse(client, branchName)

				if err != nil {
					color.Red(err.Error())
					return err
				}

				cmd = "open"
				args = []string{url}

				_, _ = exec.Command(cmd, args...).Output()

				return nil
			},
		},
		{
			Name:  "comment",
			Usage: "Comment on the ticket for the branch you are working on",
			Action: func(c *cli.Context) error {
				bytes, err := ioutil.ReadAll(os.Stdin)

				if err != nil {
					color.Red("Could not read stdin")
					return err
				}

				comment := string(bytes)

				cmd := "git"
				args := []string{"rev-parse", "--abbrev-ref", "HEAD"}

				out, err := exec.Command(cmd, args...).Output()

				if err != nil {
					message := fmt.Sprintf("Could not post message: %s. This can happen on empty repos", err)
					color.Red(message)
					return err
				}

				branchName := string(out)

				client, err := gong.NewAuthenticatedClient()

				if err != nil {
					color.Red(err.Error())
					return err
				}

				err = gong.Comment(client, branchName, comment)

				if err != nil {
					color.Red(err.Error())
					return err
				}

				return nil

			},
		},
		{
			Name:  "prepare-commit-message",
			Usage: "This is a prepare-commit-message hook for git",
			Action: func(c *cli.Context) error {
				client, err := gong.NewAuthenticatedClient()

				if err != nil {
					color.Red("Problem with starting the issue")
					return err
				}

				cmd := "git"
				args := []string{"rev-parse", "--abbrev-ref", "HEAD"}

				out, err := exec.Command(cmd, args...).Output()

				if err != nil {
					return err
				}

				branchName := string(out)
				bytes, err := ioutil.ReadAll(os.Stdin)

				if err != nil {
					color.Red("Could not read stdin")
					return err
				}

				commitMessage := string(bytes)
				newCommitMessge := gong.PrepareCommitMessage(client, branchName, commitMessage)

				fmt.Println(newCommitMessge)

				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Open the browser on the create ticket page",
			Action: func(c *cli.Context) error {
				client, err := gong.NewAuthenticatedClient()

				if err != nil {
					color.Red("Problem starting a client")
					return err
				}

				url, err := gong.Create(client)

				if err != nil {
					color.Red(err.Error())
					return err
				}

				cmd := "open"
				args := []string{url}

				_, _ = exec.Command(cmd, args...).Output()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
