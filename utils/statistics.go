package utils

import "github.com/fatih/color"

// Scalar statistic: keep track of total(or count) and average
type ScalarStatistic struct {
	total_ int
	count_ int
	name_  string
}

func (s *ScalarStatistic) Set(n string) {
	s.name_ = n
}
func (s *ScalarStatistic) Incr() {
	s.total_++
	s.count_++
}

func (s *ScalarStatistic) Add(x int) {
	s.total_ += x
	s.count_++
}
func (s *ScalarStatistic) Resume(out *color.Color) {
	out.Printf("\n %s: %d", s.name_, s.total_)
}
func (s *ScalarStatistic) ResumeAv(out *color.Color) {
	out.Printf("\n %s av: %f", s.name_, float64(s.total_)/float64(s.count_))
}

// Range statistic: keep track of min-max values
type RangeStatistic struct {
	total_ int
	count_ int
	min_   int
	max_   int
	name_  string
}

func (s *RangeStatistic) Set(n string) {
	s.name_ = n
}

func (s *RangeStatistic) Add(x int) {
	s.total_ += x
	s.count_++

	if s.count_ == 1 {
		s.max_ = x
		s.min_ = x
	} else {
		if x < s.min_ {
			s.min_ = x
		}
		if s.max_ < x {
			s.max_ = x
		}
	}
}

func (s *RangeStatistic) ResumeRange(out *color.Color) {
	out.Printf("\n %s: min=%d, max=%d", s.name_, s.min_, s.max_)
}

func (s *RangeStatistic) ResumeRangeAv(out *color.Color) {
	out.Printf("\n %s: min=%d, max=%d av=%f", s.name_, s.min_, s.max_, float64(s.total_)/float64(s.count_))
}