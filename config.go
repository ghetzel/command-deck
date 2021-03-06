package main

import (
	"os"
	"strings"

	"github.com/ghetzel/go-stockutil/executil"
	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghodss/yaml"
)

var SegmentSeparator = func() string {
	if os.Getenv(`XTERM_VERSION`) != `` {
		return ``
	} else {
		return "\uE0B0"
	}
}()

var SegmentPadding = 1

type Configuration struct {
	Segments          []*Segment `json:"segments"`
	Separator         string     `json:"separator"`
	Padding           int        `json:"padding"`
	TrailingSeparator string     `json:"trailer"`
	DisableEscape     bool       `json:"no_escape"`
	DisableColor      bool       `json:"no_colors"`
	PreCommands       []string   `json:"precommands"`
	PostCommands      []string   `json:"postcommands"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Segments:          make([]*Segment, 0),
		Separator:         SegmentSeparator,
		Padding:           SegmentPadding,
		TrailingSeparator: SegmentSeparator,
		DisableEscape:     true,
		DisableColor:      true,
	}
}

func LoadConfiguration(filename string) (*Configuration, error) {
	if data, err := fileutil.ReadAll(fileutil.MustExpandUser(filename)); err == nil {
		config := Configuration{
			Separator:         SegmentSeparator,
			Padding:           SegmentPadding,
			TrailingSeparator: SegmentSeparator,
		}

		// parse the config
		if err := yaml.Unmarshal(data, &config); err == nil {
			// since the segments were populated from data, we need to run through and
			// setup their unexported variables (pointers to config and their previous siblings)
			for i, seg := range config.Segments {
				seg.config = &config

				if seg.Padding == 0 && config.Padding != 0 {
					seg.Padding = config.Padding
				}

				if i > 0 {
					seg.prev = config.Segments[i-1]
				}
			}

			return &config, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// Create a new segment and append it to the current segments list.
func (self *Configuration) Append(expr interface{}, fg interface{}, bg interface{}) {
	var prev *Segment

	if l := len(self.Segments); l > 0 {
		prev = self.Segments[l-1]
	}

	self.Segments = append(self.Segments, &Segment{
		Expression: expr,
		FG:         fg,
		BG:         bg,
		config:     self,
		prev:       prev,
	})
}

// Append a terminator segment that resets everything.
func (self *Configuration) Close() error {
	if len(self.Segments) > 0 {
		self.Segments = append(self.Segments, &Segment{
			FG:         `default`,
			BG:         `default`,
			config:     self,
			terminator: true,
			prev:       self.Segments[len(self.Segments)-1],
		})
	}

	return nil
}

// Get all segment strings as one big line suitable for printing.
func (self *Configuration) String() string {
	var out string

	for i, cmdline := range self.PreCommands {
		if err := execAndEval(cmdline); err != nil {
			log.Warningf("bad precommand %d: %v", i+1, err)
		}
	}

	for _, seg := range self.Segments {
		if seg.enabled() {
			out += seg.String()
		}
	}

	for i, cmdline := range self.PostCommands {
		if err := execAndEval(cmdline); err != nil {
			log.Warningf("bad postcommand %d: %v", i+1, err)
		}
	}

	return out
}

func execAndEval(cmdline string) error {
	if out, err := executil.ShellCommand(cmdline).Output(); err == nil {
		for _, line := range stringutil.SplitLines(out, "\n") {
			line = strings.TrimSpace(line)

			var k, v = stringutil.SplitPair(line, `=`)

			if k != `` && v != `` {
				os.Setenv(k, v)
			}
		}

		return nil
	} else {
		return err
	}
}
