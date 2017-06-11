package util

import "strings"

type AdjustPad struct {
	Pad       string
	Count     int
	Increment int
}

func NewPad(pad string, count int) *AdjustPad {
	return &AdjustPad{
		Pad:       pad,
		Count:     count,
		Increment: 3,
	}
}

func (p *AdjustPad) LeftPad(s string) string {
	sLen := len(s)
	if sLen > p.Count {
		p.Count = sLen + p.Increment
	}

	return strings.Repeat(p.Pad, p.Count-sLen) + s
}
