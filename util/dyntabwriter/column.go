package dyntabwriter

import (
	"bytes"
)

type Column struct {
	Pad       byte
	Width     int
	Increment int
}

func NewColumn(pad byte, width int) *Column {
	return &Column{
		Pad:       pad,
		Width:     width,
		Increment: 3,
	}
}

func (c *Column) LeftPad(b []byte) []byte {
	bLen := len(b)
	if bLen > c.Width {
		c.Width = bLen + c.Increment
	}
	// TODO(leeola): copy the src slice to a new slice
	return append(bytes.Repeat([]byte{c.Pad}, c.Width-bLen), b...)
}

func (c *Column) RightPad(src []byte) []byte {
	srcLen := len(src)
	visibleSrcLen := NoColorLen(src)
	if visibleSrcLen > c.Width {
		c.Width = visibleSrcLen + c.Increment
	}

	padWidth := c.Width - visibleSrcLen
	dstLen := srcLen + padWidth
	dst := make([]byte, dstLen)
	copy(dst, src)
	for i := srcLen; i < dstLen; i++ {
		dst[i] = c.Pad
	}
	return dst
}

// NoColorLen checks the length of b while ignoring terminal escapes.
func NoColorLen(b []byte) int {
	var (
		l      int
		inRune bool
	)
	// TODO(leeola): convert to rune iteration?
	s := string(b)
	for _, c := range s {
		if c == 27 {
			inRune = true
		}
		if inRune {
			if c == 109 {
				inRune = false
			}
			continue
		}
		l++
	}
	return l
}
