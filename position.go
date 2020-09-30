package main

import (
	"fmt"
	"strconv"
	"strings"
)

var atksKnights [64]bitBoard
var atksKings [64]bitBoard

// initialize all possible knight attacks
func initAtksKnights() {
	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		// NNE  2,1
		if rk+2 < 8 && fl+1 < 8 {
			to := (rk+2)*8 + fl + 1
			toBB.set(to)
		}

		// ENE  1,2
		if rk+1 < 8 && fl+2 < 8 {
			to := (rk+1)*8 + fl + 2
			toBB.set(to)
		}

		// ESE  -1,2
		if rk-1 >= 0 && fl+2 < 8 {
			to := (rk-1)*8 + fl + 2
			toBB.set(to)
		}

		// SSE  -2,+1
		if rk-2 >= 0 && fl+1 < 8 {
			to := (rk-2)*8 + fl + 1
			toBB.set(to)
		}

		// NNW  2,-1
		if rk+2 < 8 && fl-1 >= 0 {
			to := (rk+2)*8 + fl - 1
			toBB.set(to)
		}

		// WNW  1,-2
		if rk+1 < 8 && fl-2 >= 0 {
			to := (rk+1)*8 + fl - 2
			toBB.set(to)
		}

		// WSW  -1,-2
		if rk-1 >= 0 && fl-2 >= 0 {
			to := (rk-1)*8 + fl - 2
			toBB.set(to)
		}

		// SSW  -2,-1
		if rk-2 >= 0 && fl-1 >= 0 {
			to := (rk-2)*8 + fl - 1
			toBB.set(to)
		}
		atksKnights[fr] = toBB
	}

}

// initialize all possible King attacks
func initAtksKings() {

	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		//N 1,0
		if rk+1 < 8 {
			to := (rk+1)*8 + fl
			toBB.set(to)
		}

		//NE 1,1
		if rk+1 < 8 && fl+1 < 8 {
			to := (rk+1)*8 + fl + 1
			toBB.set(to)
		}

		//E   0,1
		if fl+1 < 8 {
			to := (rk)*8 + fl + 1
			toBB.set(to)
		}

		//SE -1,1
		if rk-1 >= 0 && fl+1 < 8 {
			to := (rk-1)*8 + fl + 1
			toBB.set(to)
		}

		//S  -1,0
		if rk-1 >= 0 {
			to := (rk-1)*8 + fl
			toBB.set(to)
		}

		//SW -1,-1
		if rk-1 >= 0 && fl-1 >= 0 {
			to := (rk-1)*8 + fl - 1
			toBB.set(to)
		}

		//W   0,-1
		if fl-1 >= 0 {
			to := (rk)*8 + fl - 1
			toBB.set(to)
		}

		//NW  1,-1
		if rk+1 < 8 && fl-1 >= 0 {
			to := (rk+1)*8 + fl - 1
			toBB.set(to)
		}
		atksKings[fr] = toBB
	}
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
	b.castlings = 0
	for ix := 0; ix < nP12; ix++ {
		b.count[ix] = 0
	}

	// bitboards
	b.wbBB[WHITE], b.wbBB[BLACK] = 0, 0
	for ix := 0; ix < nP12; ix++ {
		b.pieceBB[ix] = 0
	}
}

// make a move
func (b *boardStruct) move(mv move) bool {
	newEp := 0
	// we assume that the move is legally correct (except for inCheck())
	fr := mv.fr()
	to := mv.to()
	pr := mv.pr()
	pc := b.sq[fr]

	switch {
	case p12 == wK:
		b.castlings.off(shortW | longW)
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
		b.castlings.off(shortB | longB)
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
		if (to - fr) == 16 {
			newEp = fr + 8
		} else if (to - fr) == 7 { // must be ep
			b.setSq(empty, to-8)
		} else if (to - fr) == 9 { // must be ep
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

	b.stm = b.stm ^ 0x1
	if b.isAttacked(b.King[b.stm^0x1], b.stm) {
		b.unmove(mv)
		return false
	}

	return true
}

// put the pieces back on the right squares
func (b *boardStruct) unmove(mv move) {
	b.ep = mv.ep(b.stm.opp())
	b.castlings = mv.castl()
	pc := int(mv.pc())
	fr := int(mv.fr())
	to := int(mv.to())

	b.setSq(mv.cp(), to)
	b.setSq(pc, fr)

	if pc2pt(pc) == Pawn {
		if to == b.ep && b.ep != 0 { // ep move
			b.setSq(empty, to)
			switch to - fr {
			case NW, NE:
				b.setSq(bP, to-N)
			case SW, SE:
				b.setSq(wP, to-S)
			}
		}
	} else if pc2pt(pc) == King {
		sd := pcColor(pc)
		if fr-to == 2 { // long castling
			b.setSq(castl[sd].rook, int(castl[sd].rookL))
			b.setSq(empty, fr-1)
		} else if fr-to == -2 { // short castling
			b.setSq(castl[sd].rook, int(castl[sd].rookSh))
			b.setSq(empty, fr+1)
		}
	}
	b.key = ^b.key
	b.stm = b.stm ^ 0x1
}

// check if short castlings is legal
func (b *boardStruct) isShortOk(sd color) bool {
	if !b.shortFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castl[sd].pawnsSh&b.pieceBB[Pawn]&b.wbBB[opp] != 0 { // stopped by pawns?
		return false
	}
	if castl[sd].pawnsSh&b.pieceBB[King]&b.wbBB[opp] != 0 { // stopped by king?
		return false
	}
	if castl[sd].knightsSh&b.pieceBB[Knight]&b.wbBB[opp] != 0 { // stopped by Knights?
		return false
	}

	// sliding to e1/e8	//NOTE: Maybe not needed during search because we know if we are in check
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// slidings to f1/f8
	if (mBishopTab[sq+1].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+1].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// slidings to g1/g8		//NOTE: Maybe not needed because we always make isAttacked() after a move
	if (mBishopTab[sq+2].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+2].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

// check if long castlings is legal
func (b *boardStruct) isLongOk(sd color) bool {
	if !b.longFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castl[sd].pawnsL&b.pieceBB[Pawn]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].pawnsL&b.pieceBB[King]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].knightsL&b.pieceBB[Knight]&b.wbBB[opp] != 0 {
		return false
	}

	// sliding e1/e8
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding d1/d8
	if (mBishopTab[sq-1].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-1].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding c1/c8	//NOTE: Maybe not needed because we always make inCheck() before a move
	if (mBishopTab[sq-2].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-2].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

func (b *boardStruct) setSq(p12, s int) {
	// TO IMPLEMENT
	b.sq[s] = p12
	if p12 == empty {
		b.wbBB[WHITE].clr(uint(s))
		b.wbBB[BLACK].clr(uint(s))
		for p := 0; p < nPt; p++ {
			b.pieceBB[p].clr(uint(s))
		}
	}

}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}

func (b *boardStruct) genRookMoves(ml *moveList, sd color) {
	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := uint(pc2P12(Rook, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castlings)
	var mv move

	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		toBB := mRookTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {

			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castl)
			ml.add(mv)
		}

	}

}

func (b *boardStruct) genBishopMoves(ml *moveList, sd color) {

	// TODO: generate rook moves with magic bitBoards
	allRBB := b.pieceBB[Bishop] & b.wbBB[sd]
	p12 := uint(pc2P12(Bishop, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castlings)
	var mv move

	for fr := allBBB.firstOne(); fr != 64; fr = allBBB.firstOne() {
		toBB := mBishopTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {

			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castl)
			ml.add(mv)
		}

	}

}

func (b *boardStruct) genQueenMoves(ml *moveList, sd color) {
	allRBB := b.pieceBB[Queen] & b.wbBB[sd]
	p12 := uint(pc2P12(Queen, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castlings)
	var mv move

	for fr := allQBB.firstOne(); fr != 64; fr = allQBB.firstOne() {
		toBB := mBishopTab[fr].atks(b) & (^b.wbBB[sd])
		toBB |= mRookTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {

			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castl)
			ml.add(mv)
		}

	}
}

func (b *boardStruct) genKnightMoves(ml *moveList, sd color) {
	allNBB := b.pieceBB[Knight] & b.wbBB[sd]
	p12 := uint(pc2P12(Knight, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castling)
	var mv move

	for fr := allNBB.firstOne(); fr != 64; fr = allNBB.firstOne() {
		toBB := atksKnights[fr] & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castl)
			ml.add(mv)
		}

	}
}

func (b *boardStruct) genKingMoves(ml *moveList, targetBB bitBoard) {
	sd := b.stm
	// 'normal' moves
	pc := pt2pc(King, color(sd))
	ep := b.ep
	castlings := b.castlings
	var mv move

	toBB := atksKings[b.King[sd]] & targetBB
	for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
		mv.packMove(b.King[sd], to, pc, b.sq[to], empty, ep, castlings)
		ml.add(mv)
	}

	// castlings
	if b.King[sd] == castl[sd].kingPos { // NOTE: Maybe not needed. We should know that the king is there if the flags are ok
		if targetBB.test(b.King[sd] + 2) {
			// short castling
			if b.sq[castl[sd].rookSh] == castl[sd].rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
				(castl[sd].betweenSh&b.allBB()) == 0 {
				if b.isShortOk(sd) {
					mv.packMove(b.King[sd], b.King[sd]+2, b.sq[b.King[sd]], empty, empty, b.ep, b.castlings)
					ml.add(mv)
				}
			}
		}

		if targetBB.test(b.King[sd] - 2) {
			// long castling
			if b.sq[castl[sd].rookL] == castl[sd].rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
				(castl[sd].betweenL&b.allBB()) == 0 {
				if b.isLongOk(sd) {
					mv.packMove(b.King[sd], b.King[sd]-2, b.sq[b.King[sd]], empty, empty, b.ep, b.castlings)
					ml.add(mv)
				}
			}
		}
	}
}

var genPawns = [2]func(*boardStruct, *moveList, color){(*boardStruct).genWPawnMoves, (*boardStruct).genBPawnMoves}

func (b *boardStruct) genPawnMoves(ml *moveList, sd color) {
	genPawns[sd](b, ml, sd)
}

func (b *boardStruct) genWPawnMoves(ml *moveList, sd color) {
	var mv move
	wPawns = b.pieceBB[Pawn] & b.wbBB[WHITE]

	// one step
	to1step := (wPawns << N) & ^b.allBB()

	// two steps
	to2step := ((to1step & row3) << N) ^ b.allBB()

	// captures
	toCapL := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCapR := ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]

	// promotion
	prom := (to1step | toCapL | toCapR) & row8
	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := empty
			if b.sq[to] != empty {
				// promotion capture
				cp = b.sq[to]
				if toCapL.test(uint(to)) {
					fr := to - NW
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}

				if toCapR.test(uint(to)) {
					fr := to - NE
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}
			} else {
				// promotion non-capture
				fr := to - N
				mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
				ml.add(mv)
			}

		}

	}

	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)

		// ep left
		epToL := ((wPawns & ^fileA) << NW) & epBB
		if epToL != 0 {
			mv.packMove(uint(b.ep-NW), uint(b.ep), wP, bP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}

		// ep right
		epToR := ((wPawns & ^fileH) << NE) & epBB
		if epToR != 0 {
			mv.packMove(uint(b.ep-NE), uint(b.ep), wP, bP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
	}
	// one step forward
	for to := to1step.firstOne(); to != 64; to = to1step.firstOne() {
		mv.packMove(uint(to-N), uint(to), wP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}
	// two steps forward
	for to := to2step.firstOne(); to != 64; to = to1step.firstOne() {
		mv.packMove(uint(to-2*N), uint(to), wP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}

}

func (b *boardStruct) genBPawnMoves(ml *moveList, sd color) {
	// NOT IMPLEMENTED!
}

func (b *boardStruct) genAllMoves(ml *moveList, sd color) {
	b.genpawnMoves(ml, sd)
	b.genKingMoves(ml, sd)
	b.genBishopMoves(ml, sd)
	b.genRookMoves(ml, sd)
	b.genQueenMoves(ml, sd)
	b.genKingMoves(ml, sd)

	// filter all legal moves
	b.filterLegals(ml)
}

// generate all legal moves
func (b *boardStruct) filterLegals(ml *moveList) {
	for ix := len(*ml) - 1; ix >= 0; ix-- {
		mov := (*ml)[ix]
		if b.move(mov) {
			b.unmove(mov)
		} else {
			ml.remove(ix)
		}
	}
}

func (b *boardStruct) genFrMoves(p12 int, frBB bitBoard, ml *moveList) {

}

// is sq attacked by the sd color side
func (b *boardStruct) isAttacked(to int, sd color) bool {
	if isPawnAtkingSq[sd](b, to) {
		return true
	}

	if atksKnights[to]&b.pieceBB[Knight]&b.wbBB[sd] != 0 {
		return true
	}
	if atksKings[to]&b.pieceBB[King]&b.wbBB[sd] != 0 {
		return true
	}
	if (mBishopTab[to].atks(b.allBB()) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}
	if (mRookTab[to].atks(b.allBB()) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}

	return false
}

// allAttacks from color to any square, empty or not
func (b *boardStruct) attacksBB(us color) bitBoard {
	allSq := ^bitBoard(0) // all squares

	atkBB := atksKings[b.King[us]]

	atkBB |= allPawnAtksBB[us](b)

	frBB := b.pieceBB[Knight] & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= atksKnights[fr]
	}

	frBB = (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= mBishopTab[fr].atks(allSq)
	}

	frBB = (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[us]
	for fr := frBB.firstOne(); fr != 64; fr = frBB.firstOne() {
		atkBB |= mRookTab[fr].atks(allSq)
	}

	return atkBB
}

var isPawnAtkingSq = [2]func(*boardStruct, int) bool{(*boardStruct).iswPawnAtkingSq, (*boardStruct).isbPawnAtkingSq}
var allPawnAtksBB = [2]func(*boardStruct) bitBoard{(*boardStruct).wPawnAtksBB, (*boardStruct).bPawnAtksBB}
var pawnAtksFr = [2]func(*boardStruct, int) bitBoard{(*boardStruct).wPawnAtksFr, (*boardStruct).bPawnAtksFr}
var pawnAtkers = [2]func(*boardStruct) bitBoard{(*boardStruct).wPawnAtkers, (*boardStruct).bPawnAtkers}

// Returns true or false if to-sq is attacked by white pawn
func (b *boardStruct) iswPawnAtkingSq(to int) bool {
	sqBB := bitBoard(1) << uint(to)

	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// Attacks left and right
	toCap := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCap |= ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]
	return (toCap & sqBB) != 0
}

// Returns true or false if to-sq is attacked by white pawn
func (b *boardStruct) isbPawnAtkingSq(to int) bool {
	sqBB := bitBoard(1) << uint(to)

	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// Attacks left and right
	toCap := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCap |= ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]

	return (toCap & sqBB) != 0
}

// returns all w pawns that attacka black pieces
func (b *boardStruct) wPawnAtkers() bitBoard {

	BB := b.wbBB[BLACK] // all their pieces
	// pretend that all their pieces are pawns
	// Get pawn Attacks left and right from their pieces into our pawns that now are all our pwan attackers
	ourPawnAttackers := ((BB & ^fileA) >> (-SW)) & b.wbBB[WHITE] & b.pieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^fileH) >> (-SE)) & b.wbBB[WHITE] & b.pieceBB[Pawn]

	return ourPawnAttackers
}

// returns all bl pawns that attacks white pieces
func (b *boardStruct) bPawnAtkers() bitBoard {

	BB := b.wbBB[WHITE] // all their pieces
	// pretend that all their pieces are pawns
	// Get pawn Attacks left and right from their pieces into our pawns that now are all our pwan attackers
	ourPawnAttackers := ((BB & ^fileA) << NW) & b.wbBB[BLACK] & b.pieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^fileH) << NE) & b.wbBB[BLACK] & b.pieceBB[Pawn]

	return ourPawnAttackers
}

// returns captures from fr-sq
func (b *boardStruct) wPawnAtksFr(fr int) bitBoard {
	frBB := bitBoard(1) << uint(fr)

	// Attacks left and right
	toCap := ((frBB & ^fileA) << NW) & b.wbBB[BLACK]
	toCap |= ((frBB & ^fileH) << NE) & b.wbBB[BLACK]
	return toCap
}

// returns captures from fr-sq
func (b *boardStruct) bPawnAtksFr(fr int) bitBoard {
	frBB := bitBoard(1) << uint(fr)

	// Attacks left and right
	toCap := ((frBB & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCap |= ((frBB & ^fileH) >> (-SE)) & b.wbBB[WHITE]

	return toCap
}

// returns bitBoard with all attacks, empty or not, from all white Pawns
func (b *boardStruct) wPawnAtksBB() bitBoard {
	frBB := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// Attacks left and right
	toCap := ((frBB & ^fileA) << NW)
	toCap |= ((frBB & ^fileH) << NE)
	return toCap
}

// returns bitBoard with all attacks, empty or not, from all white Pawns
func (b *boardStruct) bPawnAtksBB() bitBoard {
	frBB := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// Attacks left and right
	toCap := ((frBB & ^fileA) << NW)
	toCap |= ((frBB & ^fileH) << NE)
	return toCap
}

//////////////////////////////////// my own commands - NOT UCI /////////////////////////////////////

func (b *boardStruct) Print() {
	fmt.Println()
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}
	key, fullKey := b.key, b.fullKey()
	index := fullKey & uint64(trans.mask)
	lock := trans.lock(fullKey)
	fmt.Printf("%v to move; ep: %v  castling:%v fullKey=%x key=%x index=%x lock=%x \n", txtStm, txtEp, b.castlings.String(), fullKey, key, index, lock)

	fmt.Println("  +------+------+------+------+------+------+------+------+")
	for lines := 8; lines > 0; lines-- {
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Printf("%v |", lines)
		for ix := (lines - 1) * 8; ix < lines*8; ix++ {
			if b.sq[ix] == bP {
				fmt.Printf("   o  |")
			} else {
				fmt.Printf("   %v  |", pc2Fen(b.sq[ix]))
			}
		}
		fmt.Println()
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Println("  +------+------+------+------+------+------+------+------+")
	}

	fmt.Printf("       A      B      C      D      E      F      G      H\n")
}

func (b *boardStruct) printAllBB() {
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}
	fmt.Printf("%v to move; ep: %v   castling:%v\n", txtStm, txtEp, b.castlings.String())

	fmt.Println("white pieces")
	fmt.Println(b.wbBB[WHITE].Stringln())
	fmt.Println("black pieces")
	fmt.Println(b.wbBB[BLACK].Stringln())

	fmt.Println("wP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wK")
	fmt.Println((b.pieceBB[King] & b.wbBB[WHITE]).Stringln())

	fmt.Println("bP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bK")
	fmt.Println((b.pieceBB[King] & b.wbBB[BLACK]).Stringln())
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
	board.castlings = 0
	if len(remaining) > 1 {
		// TO IMPLEMENT
		board.castlings = parseCastling(remaining[1])
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
	row1     = bitBoard(0x00000000000000FF)
	row2     = bitBoard(0x000000000000FF00)
	row3     = bitBoard(0x0000000000FF0000)
	row4     = bitBoard(0x00000000FF000000)
	row5     = bitBoard(0x000000FF00000000)
	row6     = bitBoard(0x0000FF0000000000)
	row7     = bitBoard(0x00FF000000000000)
	row8     = bitBoard(0xFF00000000000000)
	fileA    = bitBoard(0x0101010101010101)
	fileB    = bitBoard(0x0202020202020202)
	fileG    = bitBoard(0x0404040404040404)
	fileH    = bitBoard(0x0808080808080808)
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

////////////////////////////////////// TODO: remove this //////////////////////////////////////

func (b *boardStruct) genSimplenRookMoves(ml *moveList, sd color) {

	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := uint(pc2P12(Rook, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castlings)
	var mv move
	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		rk := fr / 8
		fl := fr % 8

		// N
		for r := rk + 1; r < 8; r++ {
			to := uint(r*8 + fl)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		// S
		for r := rk - 1; r < 8; r-- {
			to := uint(r*8 + fl)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		// E
		for f := fl + 1; f < 8; f++ {
			to := uint(rk*8 + f)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		// W
		for f := fl - 1; f < 8; f-- {
			to := uint(rk*8 + f)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
	}

}
