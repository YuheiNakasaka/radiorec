package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YuheiNakasaka/radiorec/internal/recorder/ag"
	"github.com/YuheiNakasaka/radiorec/internal/recorder/radiko"
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
				cli.StringFlag{
					Name:  "station, s",
					Usage: "set radio station type(ag, radiko)",
				},
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
				switch c.String("station") {
				case "ag":
					return ag.Start(c.Int("id"), c.Int("time"))
				case "radiko":
					return radiko.Start(c.Int("id"), c.Int("time"))
				default:
					return fmt.Errorf("radio station not found(e.g -s ag)")
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
