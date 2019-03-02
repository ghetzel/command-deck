package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var rxShellExpr = regexp.MustCompile(`(\$\{!\s*(?P<shell>[^!]+)\s*!\})`) // ${! SHELL !}

func mustx(cmd string, args ...string) string {
	if out, err := x(cmd, args...); err == nil {
		return out
	} else {
		log.Warningf("exec %v %s: %v", cmd, strings.Join(args, ` `), err)
		return ``
	}
}

func x(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin

	out, err := command.Output()
	return strings.TrimSpace(string(out)), err
}

func TerminalSize() (int, int) {
	height, width := stringutil.SplitPair(mustx(`stty`, `size`), ` `)

	return int(typeutil.Int(width)), int(typeutil.Int(height))
}

func ExpandShell(in string) string {
	for {
		if match := rxutil.Match(rxShellExpr, in); match != nil {
			shell := match.Group(`shell`)

			if shell != `` {
				in = match.ReplaceGroup(1, mustx(`bash`, `-c`, shell))
			}
		} else {
			break
		}
	}

	return in
}
