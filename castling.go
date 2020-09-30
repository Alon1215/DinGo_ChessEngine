package main

import (
	"strings"
)

type castlings uint

const (
	shortW = uint(0x1) // white can casle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can casle short
	longB  = uint(0x8) // black can castle long
)

type castlOptions struct {
	short                                uint // flag
	long                                 uint // flag
	rook                                 int  // rook pc (wR/bR)
	kingPos                              int  // king pos
	rookSh                               uint // rook pos short
	rookL                                uint // rook pos long
	betweenSh                            bitBoard
	betweenL                             bitBoard
	pawnsSh, pawnsL, knightsSh, knightsL bitBoard
}

var castl = [2]castlOptions{
	{shortW, longW, wR, E1, H1, A1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{shortB, longB, bR, E8, H8, A8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
}

// castling privileges
func (c castlings) flags(sd color) bool {
	return c.shortFlag(sd) || c.longFlag(sd)
}
func (c castlings) shortFlag(sd color) bool {
	return (castl[sd].short & uint(c)) != 0
}
func (c castlings) longFlag(sd color) bool {
	return (castl[sd].long & uint(c)) != 0
}

func (c *castlings) on(val uint) {
	(*c) |= castlings(val)
}

func (c *castlings) off(val uint) {
	(*c) &= castlings(^val)
}

func (c castlings) String() string {
	flags := ""

	if uint(c)&shortW != 0 {
		flags += "K"
	}
	if uint(c)&longB != 0 {
		flags += "Q"
	}
	if uint(c)&shortB != 0 {
		flags += "k"
	}
	if uint(c)&longB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}

// parse castling rights in fenstring
func parseCastling(fenCastl string) castlings {
	c := uint(0)
	if fenCastl == "-" {
		return castlings(0)
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

	return castlings(c)
}
