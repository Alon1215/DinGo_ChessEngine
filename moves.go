package main

//directions
const (
	E          = 1
	W          = -1
	N          = 8
	S          = -8
	NW         = 7
	NE         = 9
	SW         = -NE
	SE         = -NW
	toShift    = 6
	p12Shift   = 6 + 6
	cpShift    = 4 + 6 + 6
	prShift    = 4 + 4 + 6 + 6
	epShift    = 4 + 4 + 4 + 6 + 6
	castlShift = 6 + 4 + 4 + 4 + 6 + 6
)

var pieceRules [nPt][]int //not pawns

type move uint64

func (m *move) packMove(fr, to, p12, cp, pr, ep int, castl castlings) {
	// 6 bits (fr), 6 bita (to), 4 bits (p12), 4 bits (cp), 4 bits (pr), 6 bits (ep), 4 bits (castl), x bits value
	*m = move(fr | (to << toShift) | (p12 | p12Shift) | (cp << cpShift) | (pr << prShift) | (ep << epShift) | uint(castl<<castlShift))
}

func init() {
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], W)
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], S)
}

func (m move) fr() int {
	return int(m & frMask)
}
func (m move) to() int {
	return int(m&toMask) >> toShift
}

func (m move) pc() int {
	return int(m&pcMask) >> pcShift
}
func (m move) cp() int {
	return int(m&cpMask) >> cpShift
}
func (m move) pr() int {
	return int(m&prMask) >> prShift
}

func (m move) ep(sd color) int {
	// sd is the side that can capture
	file := int(m&epMask) >> epShift
	if file == 0 {
		return 0 // no ep
	}

	// there is an ep sq
	rank := 5
	if sd == BLACK {
		rank = 2
	}

	return rank*8 + file - 1

}

func (m move) castl() castlings {
	return castlings(m&castlMask) >> castlShift
}

//move without eval
func (m move) onlyMv() move {
	return m & move(^evalMask)
}

type moveList []move

func (ml *moveList) new(size int) {
	*ml = make(moveList, 0, size)
}

func (ml *moveList) clear() {
	*ml = (*ml)[:0]
}

func (ml *moveList) add(mv move) {
	*ml = append(*ml, mv)
}

func (ml *moveList) remove(ix int) {
	if len(*ml) > ix && ix >= 0 {
		*ml = append((*ml)[:ix], (*ml)[ix+1:]...)
	}
}

// Sort is sorting the moves in the Score/Move list according to the score per move
func (ml *moveList) sort() {
	bSwap := true
	for bSwap {
		bSwap = false
		for i := 0; i < len(*ml)-1; i++ {
			if (*ml)[i+1].eval() > (*ml)[i].eval() {
				(*ml)[i], (*ml)[i+1] = (*ml)[i+1], (*ml)[i]
				bSwap = true
			}
		}
	}
}

func (ml moveList) String() string {
	theString := ""
	for _, mv := range ml {
		theString += mv.String() + " "
	}
	return theString
}
