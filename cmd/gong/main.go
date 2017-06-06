package main

import (
	"fmt"
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
					fmt.Println(err)
					return nil
				}

				return gong.Login(client)
			},
		},
	}
	app.Run(os.Args)
}
