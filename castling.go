package main

import (
	"strings"
)

type castling unit

const (
	shortW = unit(0x1) // white can casle short
	longW  = unit(0x2) // white can castle long
	shortB = unit(0x4) // black can casle short
	longB  = unit(0x8) // black can castle long
)

func parseCastling(fenCastl string) castling {
	c := uint(0)
	if fenCastl == "-" {
		return castling(0)
	}

	if strings.Index(fenCastl, "K") >= 0 {
		c |= shortW
	}

	if strings.Index(fenCastl, "Q") >= 0 {
		c |= longW
	}

	if strings.Index(fenCastl, "k") >= 0 {
		c |= shortB
	}

	if strings.Index(fenCastl, "q") >= 0 {
		c |= longB
	}

	return castling(c)
}
