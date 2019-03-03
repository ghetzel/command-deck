package main

import (
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/typeutil"
)

type Segment struct {
	Name              string      `json:"name"`
	Disable           bool        `json:"disable"`
	Expression        interface{} `json:"expr"`
	FG                interface{} `json:"fg"`
	BG                interface{} `json:"bg"`
	Separator         string      `json:"separator,omitempty"`
	Padding           int         `json:"padding,omitempty"`
	Timeout           string      `json:"timeout,omitempty"`
	NotIf             string      `json:"except,omitempty"`
	OnlyIf            string      `json:"if,omitempty"`
	ReverseJoinColors bool        `json:"reverse,omitempty"`
	prev              *Segment
	config            *Configuration
	terminator        bool
}

func (self *Segment) previous() *Segment {
	if self.prev != nil {
		prev := self.prev

		for !prev.enabled() {
			prev = prev.prev

			if prev == nil {
				return nil
			}
		}

		return prev
	}

	return nil
}

func (self *Segment) enabled() bool {
	if self.terminator {
		return true
	}

	if self.Disable {
		return false
	}

	if self.NotIf != `` {
		if xerr(self.NotIf) == nil {
			return true
		}
	}

	if self.OnlyIf != `` {
		if xerr(self.OnlyIf) != nil {
			return false
		}
	}

	if typeutil.IsZero(self.Expression) {
		return false
	}

	return true
}

func (self *Segment) Foreground() string {
	if typeutil.Int(self.FG) < 0 {
		return `default`
	} else {
		return ExpandShell(typeutil.String(self.FG), self.Timeout)
	}
}

func (self *Segment) Background() string {
	if typeutil.Int(self.BG) < 0 {
		return `default`
	} else {
		return ExpandShell(typeutil.String(self.BG), self.Timeout)
	}
}

func (self *Segment) Sep() string {
	if self.terminator && self.config != nil && self.config.TrailingSeparator != `` {
		return self.config.TrailingSeparator
	} else if self.Separator != `` {
		return self.Separator
	} else if self.config != nil {
		return self.config.Separator
	} else {
		return SegmentSeparator
	}
}

func (self *Segment) Pad() int {
	if self.Padding > 0 {
		return self.Padding
	} else if self.config != nil {
		return self.config.Padding
	} else {
		return 0
	}
}

func (self *Segment) String() string {
	out := ``
	fg := self.Foreground()
	bg := self.Background()

	if self.enabled() {
		if prev := self.previous(); prev != nil {
			if self.ReverseJoinColors {
				out += "${" + bg + ":" + prev.Background() + "}"
			} else {
				out += "${" + prev.Background() + ":" + bg + "}"
			}

			out += self.Sep()
		}

		if self.terminator {
			out += "${reset} "
		} else {
			expr := typeutil.String(self.Expression)

			if expr != `` {
				out += "${" + fg + ":" + bg + "}"
				out += strings.Repeat(` `, self.Pad())
				out += expr
				out += strings.Repeat(` `, self.Pad())
			}

			out += "${reset}"
			out = ExpandShell(out, self.Timeout)
		}
	}

	if fileutil.IsTerminal() {
		out = log.CSprintf(out)
	} else {
		out = log.TermSprintf(out)
	}

	return out
}
