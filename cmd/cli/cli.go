package main

import (
	"log"
	"os"

	"github.com/YuheiNakasaka/radiorec/internal/recorder/ag"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rediorec"
	app.Version = "0.0.1"
	app.Usage = "A cli application to record specific radio programs"

	app.Commands = []cli.Command{
		{
			Name:    "record",
			Aliases: []string{"r"},
			Usage:   "Record radio",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "id, i",
					Usage: "set program ID",
				},
				cli.IntFlag{
					Name:  "time, t",
					Usage: "set airtime",
				},
			},
			Action: func(c *cli.Context) error {
				return ag.Start(c.Int("id"), c.Int("time"))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
