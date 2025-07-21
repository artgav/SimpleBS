package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"SimpleBS/localcluster"
)

func runCLI() {
	app := &cli.App{
		Name:  "SimpleBS",
		Usage: "Manage local volumes via SimpleBS",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "folder",
				Usage:    "Path to volume storage folder",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a volume",
				Action: func(c *cli.Context) error {
					v := localcluster.NewVendor(c.String("folder"))
					size, _ := strconv.ParseInt(c.Args().Get(1), 10, 64)
					return v.CreateVolume(c.Args().Get(0), size)
				},
			},
			{
				Name:  "info",
				Usage: "Get volume info",
				Action: func(c *cli.Context) error {
					v := localcluster.NewVendor(c.String("folder"))
					info, err := v.GetVolumeInfo(c.Args().Get(0))
					if err != nil {
						return err
					}
					fmt.Println(info)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a volume",
				Action: func(c *cli.Context) error {
					v := localcluster.NewVendor(c.String("folder"))
					return v.DeleteVolume(c.Args().Get(0))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
