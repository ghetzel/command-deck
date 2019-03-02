package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

const (
	BackgroundChar = '-'
)

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
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Action = func(c *cli.Context) {
		segments := new(Segments)

		if retcode := typeutil.Int(c.Args().First()); retcode != 0 {
			segments.Append(retcode, 15, 1)
		}

		// append time
		segments.Append(`${! date +%H:%M:%S !}`, 252, 240)

		// append user@host
		segments.Append(`${! whoami !}@${! hostname !}`, 15, 4)

		// append pwd
		segments.Append(`${! pwd | sed "s|$HOME|~|" !}`, 15, 64)

		// append git status (if applicable)
		if fileutil.IsNonemptyFile(`.git/HEAD`) {
			vcscolor := 24

			var ref string

			if line, err := fileutil.ReadFirstLine(`.git/HEAD`); err == nil {
				_, ref = stringutil.SplitPair(line, `: `)
				ref = strings.TrimPrefix(ref, `refs/heads/`)
			} else {
				ref = `!ERR!`
			}

			gserr := make(chan error)

			go func() {
				_, err := x(`git`, `diff-index`, `--quiet`, `HEAD`)
				gserr <- err
			}()

			select {
			case err := <-gserr:
				if err != nil {
					vcscolor = 9 // dirty
				}
			case <-time.After(250 * time.Millisecond):
				vcscolor = 240 // unknown
			}

			segments.Append(ref, 15, vcscolor)
		}

		segments.Close()
		fmt.Print(segments.String())
	}

	app.Run(os.Args)
}
