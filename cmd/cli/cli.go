package main

import (
	"fmt"
	"log"
	"os"

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
				case 2: // radiko
					recorder := radiko.Radiko{}
					return recorder.Start(program.ID, program.Airtime, c.String("storage"))
				default:
					return fmt.Errorf("the program not found: %v", c.Int("id"))
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
