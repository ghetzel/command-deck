package main

import (
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var SegmentSeparator = "\uE0B0"
var SegmentPadding = 1

type Segment struct {
	prev       *Segment
	expression interface{}
	fg         interface{}
	bg         interface{}
	sep        string
	padding    int
}

func NewSegment(prev *Segment, expression interface{}, fg interface{}, bg interface{}) *Segment {
	return &Segment{
		prev:       prev,
		expression: expression,
		fg:         fg,
		bg:         bg,
		sep:        SegmentSeparator,
		padding:    SegmentPadding,
	}
}

func (self *Segment) String() string {
	out := ``

	if typeutil.Int(self.fg) < 0 {
		self.fg = `default`
	}

	if typeutil.Int(self.bg) < 0 {
		self.bg = `default`
	}

	if self.prev != nil {
		out += "${" + typeutil.String(self.prev.bg) + ":" + typeutil.String(self.bg) + "}"
		out += self.sep
	}

	expr := typeutil.String(self.expression)
	justWhitespace := (len(strings.TrimSpace(expr)) == 0)

	if expr != `` {
		out += "${" + typeutil.String(self.fg) + ":" + typeutil.String(self.bg) + "}"

		if !justWhitespace {
			out += strings.Repeat(` `, self.padding)
		}

		out += expr

		if !justWhitespace {
			out += strings.Repeat(` `, self.padding)
		}
	}

	out += "${reset}"
	out = ExpandShell(out)

	if fileutil.IsTerminal() {
		out = log.CSprintf(out)
	} else {
		out = log.TermSprintf(out)
	}

	return out
}
