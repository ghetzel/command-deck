package main

import (
	"fmt"
	"os"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/log"
)

const defaultPS1 = `\u@\h \W \$ `

func main() {
	app := cli.NewApp()
	app.Name = `cdeck`
	app.Usage = `A fancy shell prompt.`
	app.Version = `0.0.1`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `debug`,
			EnvVar: `LOGLEVEL`,
		},
		cli.StringFlag{
			Name:  `config, c`,
			Usage: `The configuration file to load.`,
			Value: `~/.config/cdeck.yml`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Action = func(c *cli.Context) {
		if config, err := LoadConfiguration(c.String(`config`)); err == nil {
			config.Close()
			fmt.Print(config.String())
		} else {
			fmt.Print(defaultPS1)
		}
	}

	app.Run(os.Args)
}
