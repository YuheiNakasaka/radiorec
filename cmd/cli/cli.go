package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YuheiNakasaka/radiorec/internal/cron"
	"github.com/YuheiNakasaka/radiorec/internal/db"
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
				cli.IntFlag{
					Name:  "id, i",
					Usage: "set program ID",
				},
				cli.StringFlag{
					Name:  "storage, s",
					Value: "",
					Usage: "set external storage flag. Default is local. (e.g -s s3)",
				},
			},
			Action: func(c *cli.Context) error {
				mydb := &db.MyDB{}
				err := mydb.New()
				if err != nil {
					return fmt.Errorf("Failed to connect database: %v", err)
				}
				program := mydb.FindProgram(c.Int("id"))

				switch program.Station {
				case 1: // a&g
					recorder := ag.Ag{}
					return recorder.Start(program.ID, program.Airtime, c.String("storage"))
				case 2: // 文化放送
					recorder := radiko.Radiko{}
					return recorder.Start(program.ID, program.Airtime, c.String("storage"), 2)
				case 3: // hellofive
					recorder := radiko.Radiko{}
					return recorder.Start(program.ID, program.Airtime, c.String("storage"), 3)
				default:
					return fmt.Errorf("the program not found: %v", c.Int("id"))
				}
			},
		},
		{
			Name:    "cron",
			Aliases: []string{"c"},
			Usage:   "Clear old crontab and recreate new crontab.",
			Action: func(c *cli.Context) error {
				return cron.Generate()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
