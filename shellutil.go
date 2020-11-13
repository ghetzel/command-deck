package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/timeutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var rxShellExpr = regexp.MustCompile(`(\{!\s*(?P<shell>[^!]+)\s*!\})`) // {! SHELL !}
var argv []string

func shellx(shellcmd string) (string, error) {
	cmd := `bash`

	if s := os.Getenv(`SHELL`); s != `` {
		cmd = s
	}

	return x(cmd, `-c`, shellcmd)
}

func mustx(shellcmd string) string {
	if out, err := shellx(shellcmd); err == nil {
		return out
	} else {
		log.Noticef("exec %s: %v", shellcmd, err)
		return ``
	}
}

func xerr(shellcmd string) error {
	if shellcmd != `` {
		_, err := shellx(shellcmd)
		return err
	} else {
		return nil
	}
}

func x(cmd string, args ...string) (string, error) {
	var command = exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Env = os.Environ()

	for i, arg := range argv {
		command.Env = append(command.Env, fmt.Sprintf("CDECK_ARGV_%d=%s", i+1, arg))
	}

	log.Debugf("exec %s %s:", cmd, strings.Join(args, ` `))

	for _, ev := range command.Env {
		log.Debugf("  %v", ev)
	}

	out, err := command.Output()
	return strings.TrimSpace(string(out)), err
}

// func TerminalSize() (int, int) {
// 	height, width := stringutil.SplitPair(mustx(`stty`, `size`), ` `)
// 	return int(typeutil.Int(width)), int(typeutil.Int(height))
// }

func ExpandShell(in string, timeoutI interface{}) string {
	var timeout time.Duration = 1 * time.Second

	if t, err := timeutil.ParseDuration(typeutil.String(timeoutI)); err == nil && t > 0 {
		timeout = t
	}

	for {
		if match := rxutil.Match(rxShellExpr, in); match != nil {
			shell := match.Group(`shell`)

			if shell != `` {
				var outchan = make(chan string)

				go func() {
					outchan <- mustx(shell)
				}()

				select {
				case out := <-outchan:
					in = match.ReplaceGroup(1, out)

				case <-time.After(timeout):
					log.Noticef("exec %v: timeout", shell)
					in = match.ReplaceGroup(1, ``)
				}
			}
		} else {
			break
		}
	}

	return in
}
