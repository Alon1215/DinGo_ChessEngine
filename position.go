package main

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	initFenSq2Int()
}

type boardStruct struct {
	key     uint64
	sq      [64]int
	wbBB    [2]bitBoard
	pieceBB [nPt]bitBoard
	King    [2]int
	ep      int
	castlings
	stm    color
	count  [12]int
	rule50 int //set to 0 if a pawn or capt move otherwise increment
}

type color int

var board = boardStruct{}

func (b *boardStruct) allBB() bitBoard {
	return b.wbBB[0] | b.wbBB[1]
}

// Clear the board, flags, bitboard  etc

func (b *boardStruct) clear() {
	b.stm = WHITE
	b.rule50 = 0
	b.sq = [64]int{}
	for ix := A1; ix <= H8; ix++ {
		b.sq[ix] = empty
	}
	b.ep = 0
	b.castling = 0
	for ix := 0; ix < nP12; ix++ {
		b.count[ix] = 0
	}

	// bitboards
	b.wbBB[WHITE], b.wbBB[BLACK] = 0, 0
	for ix := 0; ix < nP12; ix++ {
		b.pieceBB[ix] = 0
	}
}

func (b *boardStruct) move(fr, to, pr int) bool {
	// TODO 1. move in board
	// Assumption: move is legal
	p12 := b.sq[fr]
	switch p12 {
	case p12 == wK:
		b.castling.off(shortW | longW)
		if abs(to-fr) == 2 {
			if fr == E1 {
				if to == G1 {
					b.setSq(wR, F1)
					b.setSq(empty, H1)
				} else {
					b.setSq(wR, D1)
					b.setSq(empty, A1)
				}
			}
		}
	case p12 == bK:
		b.castling.off(shortB | longB)
		if abs(to-fr) == 2 {
			if fr == E1 {
				if to == G1 {
					b.setSq(bR, F1)
					b.setSq(empty, H1)
				} else {
					b.setSq(bR, D1)
					b.setSq(empty, A1)
				}
			}
		}
	case p12 == wR:
		if fr == A1 {
			b.off(longW)
		} else if fr == H1 {
			b.off(shortW)
		}
	case p12 == bR:
		if fr == A1 {
			b.off(longB)
		} else if fr == H1 {
			b.off(shortB)
		}
	case p12 == wP && b.sq[to] == empty:
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep
			b.setSq(empty, to-8)
		} else if to-fr == 9 { // must be ep
			b.setSq(empty, to-8)
		}
		// handle ep
	case p12 == bP && b.sq[to] == empty:
		if to-fr == 16 {
			newEp = fr + 8
		} else if fr-to == 7 { // must be ep
			b.setSq(empty, to+8)
		} else if fr-to == 9 { // must be ep
			b.setSq(empty, to+8)
		}

	}
	b.ep = newEp
	b.setSq(empty, fr)
	if pr != empty {
		b.setSq(pr, to)
	} else {
		b.setSq(p12, to)
	}
	// TODO isInCheck(stm) need to be mad

	/*
		if b.isInCheck(b.stm) {
			b.stm = b.stm ^ 0x1
			return false
		}
	*/
	b.stm = b.stm ^ 0x1
	return true
}
func (b *boardStruct) setSq(p12, s int) {
	// TO IMPLEMENT
	b.sq[s] = p12
	if p12 == empty {
		b.wbBB[WHITE].clr(uint(s))
		b.wbBB[BLACK].clr(uint(s))
		for p := 0; p < nP; p++ {
			b.pieceBB[p].clr(uint(s))
		}
	}

}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}

//  parse a FEN string and setup that position
func parseFEN(FEN string) {
	fenIx := 0
	sq := 0

	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {
			char := string(FEN[fenIx])
			fenIx++
			if char == "/" {
				continue
			}

			if i, err := strconv.Atoi(char); err == nil { // numeriskt
				// fmt.Println(i, "empty from sq", sq)
				// sq += i

				for j := 0; j < i; j++ {
					board.setSq(empty, sq)
					sq++
				}
				continue
			}

			//fmt.Println(char, "on sq", sq)
			if strings.IndexAny(p12ToFen, char) == -1 {
				tell("info string invalid piece ", char, " try next one")
				continue
			}

			board.setSq(fen2Int(char), sq)

			sq++
		}
		// take care of side to move
		// take care of castling rights
		// set the 50 mpve rule
		// set number of full moves
	}
	remaining := strings.Split(trim(FEN[fenIx:]), " ")

	// stm
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			board.stm = WHITE
		} else if remaining[0] == "b" {
			board.stm = BLACK
		} else {
			r := fmt.Sprintf("%vl sq=%v; fenIx=%v", strings.Join(remaining, " "), sq, fenIx)

			tell("info string remaining", r, ";")
			tell("info string", remaining[0], " invalid stm color")
			board.stm = WHITE
		}
	}

	// castling
	board.castling = 0
	if len(remaining) > 1 {
		// TO IMPLEMENT
		board.castling = parseCastling(remaining[1])
	}

	// ep square
	board.ep = 0
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			board.ep = fenSq2Int[remaining[2]]
		}
	}

	// 50-move
	board.rule50 = 0
	if len(remaining) > 3 {
		board.parse50(remaining[3])
	}
}

// parse 50 move rule ni fenstring
func parse50(fen50 string) int {
	// TO IMPLEMENT

}

// parse and make the moves in position command from GUI

func parseMVS(mvstr string) {
	// mvs := strings.Split(mvstr, " ")
	mvs := strings.Fields(low(mvstr))

	for _, mv := range mvs {
		// fmt.Println("make move", mv)
		mv = trim(mv)
		if len(mv) < 4 {
			tell("info string ", mv, " in the position command is not a correct move")
			return
		}
		// is fr square is ok?
		fr, ok := fen2Sq2Int[mv[:2]]
		if !ok {
			tell("info string ", mv, " in the position command is not a correct fr square")
			return
		}
		p12 := board.sq[fr]
		if p12 == empty {
			tell("info string ", mv, " in the position command. fr_sq is not an empty suare")
			return
		}

		pCol := p12Color(p12)
		if pCol != board.stm {
			tell("info string ", mv, " in the position command. fr piece has the wrong color")
			return
		}

		// is to square is ok?
		to, ok := femSq2Int[mv[2:4]]
		if !ok {
			tell("info string ", mv, " in the position command is not a correct to square")
			return
		}

		// is the prom piece ok?
		pr := 0
		if len(mv) == 5 { //prom
			if !strings.ContainsAny(mv[4:5], "qrbn") {
				tell("info string promotion piece in ", mv, " in the position command is not a correct --MISS THE REST--")
				return
			}

			pr = fen2Int(mv[4:5])
			pr = pc2P12(pr, board.stm)
		}
		board.move(fr, to, pr)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// // fen2Int pieceString to p12 int
// func fen2Int(c string) int {
// 	// TO IMPLEMENT

// }

// // int2Fen pieceString to p12 int
// func int2Fen(p12 int) string {
// 	// TO IMPLEMENT
// }

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// fen2pc convert pieceString to pc int
func fen2pc(c string) int {
	for p, x := range pcFen {
		if string(x) == c {
			return p
		}
	}
	return empty
}

// pc2Fen convert pc to fenString
func pc2Fen(pc int) string {
	if pc == empty {
		return " "
	}
	return string(pcFen[pc])
}

// pc2pt returns the pt from pc
func pc2pt(pc int) int {
	return pc >> 1
}

// pcColor returns the color of a pc form
func pcColor(pc int) color {
	return color(pc & 0x1)
}

// pt2pc returns pc from pt and sd
func pt2pc(pt int, sd color) int {
	return (pt << 1) | int(sd)
}

// map fen-sq to int
var fen2Sq = make(map[string]int)

// map int-sq to fen
var sq2Fen = make(map[int]string)

// init the square map from string to int and int to string
func initFen2Sq() {
	fen2Sq["a1"] = A1
	fen2Sq["a2"] = A2
	fen2Sq["a3"] = A3
	fen2Sq["a4"] = A4
	fen2Sq["a5"] = A5
	fen2Sq["a6"] = A6
	fen2Sq["a7"] = A7
	fen2Sq["a8"] = A8

	fen2Sq["b1"] = B1
	fen2Sq["b2"] = B2
	fen2Sq["b3"] = B3
	fen2Sq["b4"] = B4
	fen2Sq["b5"] = B5
	fen2Sq["b6"] = B6
	fen2Sq["b7"] = B7
	fen2Sq["b8"] = B8

	fen2Sq["c1"] = C1
	fen2Sq["c2"] = C2
	fen2Sq["c3"] = C3
	fen2Sq["c4"] = C4
	fen2Sq["c5"] = C5
	fen2Sq["c6"] = C6
	fen2Sq["c7"] = C7
	fen2Sq["c8"] = C8

	fen2Sq["d1"] = D1
	fen2Sq["d2"] = D2
	fen2Sq["d3"] = D3
	fen2Sq["d4"] = D4
	fen2Sq["d5"] = D5
	fen2Sq["d6"] = D6
	fen2Sq["d7"] = D7
	fen2Sq["d8"] = D8

	fen2Sq["e1"] = E1
	fen2Sq["e2"] = E2
	fen2Sq["e3"] = E3
	fen2Sq["e4"] = E4
	fen2Sq["e5"] = E5
	fen2Sq["e6"] = E6
	fen2Sq["e7"] = E7
	fen2Sq["e8"] = E8

	fen2Sq["f1"] = F1
	fen2Sq["f2"] = F2
	fen2Sq["f3"] = F3
	fen2Sq["f4"] = F4
	fen2Sq["f5"] = F5
	fen2Sq["f6"] = F6
	fen2Sq["f7"] = F7
	fen2Sq["f8"] = F8

	fen2Sq["g1"] = G1
	fen2Sq["g2"] = G2
	fen2Sq["g3"] = G3
	fen2Sq["g4"] = G4
	fen2Sq["g5"] = G5
	fen2Sq["g6"] = G6
	fen2Sq["g7"] = G7
	fen2Sq["g8"] = G8

	fen2Sq["h1"] = H1
	fen2Sq["h2"] = H2
	fen2Sq["h3"] = H3
	fen2Sq["h4"] = H4
	fen2Sq["h5"] = H5
	fen2Sq["h6"] = H6
	fen2Sq["h7"] = H7
	fen2Sq["h8"] = H8

	// -------------- sq2Fen
	sq2Fen[A1] = "a1"
	sq2Fen[A2] = "a2"
	sq2Fen[A3] = "a3"
	sq2Fen[A4] = "a4"
	sq2Fen[A5] = "a5"
	sq2Fen[A6] = "a6"
	sq2Fen[A7] = "a7"
	sq2Fen[A8] = "a8"

	sq2Fen[B1] = "b1"
	sq2Fen[B2] = "b2"
	sq2Fen[B3] = "b3"
	sq2Fen[B4] = "b4"
	sq2Fen[B5] = "b5"
	sq2Fen[B6] = "b6"
	sq2Fen[B7] = "b7"
	sq2Fen[B8] = "b8"

	sq2Fen[C1] = "c1"
	sq2Fen[C2] = "c2"
	sq2Fen[C3] = "c3"
	sq2Fen[C4] = "c4"
	sq2Fen[C5] = "c5"
	sq2Fen[C6] = "c6"
	sq2Fen[C7] = "c7"
	sq2Fen[C8] = "c8"

	sq2Fen[D1] = "d1"
	sq2Fen[D2] = "d2"
	sq2Fen[D3] = "d3"
	sq2Fen[D4] = "d4"
	sq2Fen[D5] = "d5"
	sq2Fen[D6] = "d6"
	sq2Fen[D7] = "d7"
	sq2Fen[D8] = "d8"

	sq2Fen[E1] = "e1"
	sq2Fen[E2] = "e2"
	sq2Fen[E3] = "e3"
	sq2Fen[E4] = "e4"
	sq2Fen[E5] = "e5"
	sq2Fen[E6] = "e6"
	sq2Fen[E7] = "e7"
	sq2Fen[E8] = "e8"

	sq2Fen[F1] = "f1"
	sq2Fen[F2] = "f2"
	sq2Fen[F3] = "f3"
	sq2Fen[F4] = "f4"
	sq2Fen[F5] = "f5"
	sq2Fen[F6] = "f6"
	sq2Fen[F7] = "f7"
	sq2Fen[F8] = "f8"

	sq2Fen[G1] = "g1"
	sq2Fen[G2] = "g2"
	sq2Fen[G3] = "g3"
	sq2Fen[G4] = "g4"
	sq2Fen[G5] = "g5"
	sq2Fen[G6] = "g6"
	sq2Fen[G7] = "g7"
	sq2Fen[G8] = "g8"

	sq2Fen[H1] = "h1"
	sq2Fen[H2] = "h2"
	sq2Fen[H3] = "h3"
	sq2Fen[H4] = "h4"
	sq2Fen[H5] = "h5"
	sq2Fen[H6] = "h6"
	sq2Fen[H7] = "h7"
	sq2Fen[H8] = "h8"
}

// varius consts
const (
	nP12     = 12
	nPt      = 6
	WHITE    = color(0)
	BLACK    = color(1)
	startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -"
)

// piece char definitions
const (
	pc2Char  = "PNBRQK ?"
	p12ToFen = "PpNnBbRrQqKk"
)

// 6 pieces types - no color (P)
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// 12 pieces with color (P12)
const (
	wP = iota
	bP
	wN
	bN
	wB
	bB
	wR
	bR
	wQ
	bQ
	wK
	bK
	empty = 15
)

// square names
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)
