package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"simplebs/localcluster"
)

func runCLI() {
	app := &cli.App{
		Name:  "SimpleBS",
		Usage: "Manage local volumes via CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "folder",
				Usage:    "Path to volume storage folder",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "create",
				Usage:     "Create a volume: create <name> <size_bytes>",
				ArgsUsage: "<name> <size_bytes>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return cli.Exit("Usage: create <name> <size_bytes>", 1)
					}
					name := c.Args().Get(0)
					size, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
					if err != nil {
						return fmt.Errorf("invalid size: %v", err)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.CreateVolume(name, size)
				},
			},
			{
				Name:      "resize",
				Usage:     "Resize a volume: resize <name> <new_size_bytes>",
				ArgsUsage: "<name> <new_size_bytes>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return cli.Exit("Usage: resize <name> <new_size_bytes>", 1)
					}
					name := c.Args().Get(0)
					newSize, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
					if err != nil {
						return fmt.Errorf("invalid size: %v", err)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.ResizeVolume(name, newSize)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a volume: delete <name>",
				ArgsUsage: "<name>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return cli.Exit("Usage: delete <name>", 1)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.DeleteVolume(c.Args().Get(0))
				},
			},
			{
				Name:      "info",
				Usage:     "Show volume info: info <name>",
				ArgsUsage: "<name>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return cli.Exit("Usage: info <name>", 1)
					}
					v := localcluster.NewVendor(c.String("folder"))
					info, err := v.GetVolumeInfo(c.Args().Get(0))
					if err != nil {
						return err
					}
					for k, v := range info {
						fmt.Printf("%s: %v\n", k, v)
					}
					return nil
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone a volume: clone <src> <dst>",
				ArgsUsage: "<src> <dst>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return cli.Exit("Usage: clone <src> <dst>", 1)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.CloneVolume(c.Args().Get(0), c.Args().Get(1))
				},
			},
			{
				Name:      "attach",
				Usage:     "Attach volume to a server: attach <name> <server_ip>",
				ArgsUsage: "<name> <server_ip>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return cli.Exit("Usage: attach <name> <server_ip>", 1)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.AttachVolume(c.Args().Get(0), c.Args().Get(1))
				},
			},
			{
				Name:      "detach",
				Usage:     "Detach a volume: detach <name>",
				ArgsUsage: "<name>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return cli.Exit("Usage: detach <name>", 1)
					}
					v := localcluster.NewVendor(c.String("folder"))
					return v.DetachVolume(c.Args().Get(0))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
