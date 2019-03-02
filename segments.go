package main

type Segments struct {
	segments []*Segment
}

func (self *Segments) Append(expression interface{}, fg interface{}, bg interface{}) *Segment {
	var seg *Segment

	if len(self.segments) > 0 {
		seg = NewSegment(self.segments[len(self.segments)-1], expression, fg, bg)
	} else {
		seg = NewSegment(nil, expression, fg, bg)
	}

	self.segments = append(self.segments, seg)
	return seg
}

func (self *Segments) Close() error {
	self.Append(` `, -1, -1)
	return nil
}

func (self *Segments) String() string {
	var out string

	for _, seg := range self.segments {
		out += seg.String()
	}

	return out
}
