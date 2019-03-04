package main

import (
	"fmt"
	"os"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
)

const defaultPS1 = `\u@\h \W \$ `

func main() {
	app := cli.NewApp()
	app.Name = `cdeck`
	app.Usage = `A fancy shell prompt.`
	app.Version = `0.0.2`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `warning`,
			EnvVar: `CDECK_LOGLEVEL`,
		},
		cli.StringFlag{
			Name:  `config, c`,
			Usage: `The configuration file to load.`,
			Value: `~/.config/cdeck.yml`,
		},
		cli.BoolFlag{
			Name:  `eval, e`,
			Usage: `Evaluate the patterns given on the command line instead of reading from the config file.`,
		},
		cli.IntFlag{
			Name:  `padding, p`,
			Usage: `The default padding to use.`,
		},
		cli.StringFlag{
			Name:  `separator, s`,
			Usage: `The default separator string to use.`,
		},
		cli.StringFlag{
			Name:  `trailer, T`,
			Usage: `The string to use for the trailing separator.`,
		},
		cli.BoolFlag{
			Name:  `disable-term-escape, E`,
			Usage: `Disables wrapping color escape sequences with '\[' and '\]'`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Action = func(c *cli.Context) {
		var config *Configuration

		argv = c.Args()

		if c.Bool(`eval`) {
			config = NewConfiguration()

			for _, arg := range argv {
				fg, bg, expr := stringutil.SplitTriple(arg, `:`)
				config.Append(expr, fg, bg)
			}
		} else {
			if c, err := LoadConfiguration(c.String(`config`)); err == nil {
				config = c
			}
		}

		if c.IsSet(`padding`) {
			config.Padding = c.Int(`padding`)
		}

		if c.IsSet(`separator`) {
			config.Separator = c.String(`separator`)
		}

		if c.IsSet(`trailer`) {
			config.TrailingSeparator = c.String(`trailer`)
		}

		if c.IsSet(`disable-term-escape`) {
			config.DisableTermEscape = c.Bool(`disable-term-escape`)
		}

		if config != nil {
			config.Close()
			fmt.Print(config.String())
		} else {
			fmt.Print(defaultPS1)
		}
	}

	app.Run(os.Args)
}
