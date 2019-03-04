package main

import (
	"os"

	"github.com/ghetzel/go-stockutil/fileutil"
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

		if err := yaml.Unmarshal(data, &config); err == nil {
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

func (self *Configuration) String() string {
	var out string

	for _, seg := range self.Segments {
		if seg.enabled() {
			out += seg.String()
		}
	}

	return out
}
