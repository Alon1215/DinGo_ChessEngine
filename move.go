package main

//directions
const (
	E  = 1
	W  = -1
	N  = 8
	S  = -8
	NW = 7
	NE = 9
	SW = -NE
	SE = -NW
)

var pieceRules [nPt][]int //not pawns

type move uint64

func (m *move) packMove(fr, to, p12, empty, ep int, castl castling) {
	// 6 bits (fr), 6 bita (to), 4 bits (p12), 4 bits (cp), 4 bits (pr), 6 bits (ep), 4 bits (castl), x bits value

}

type moveList struct {
	mv []move
}

func init() {
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], W)
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], S)
}

func (ml *moveList) add(mv move) {
	ml.mv = append(ml.mv, mv)

}
