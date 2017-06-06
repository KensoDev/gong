package main

import (
	"github.com/fatih/color"
	"github.com/kensodev/gong"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"

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
	}
	app.Run(os.Args)
}
