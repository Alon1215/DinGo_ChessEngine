package main

// const (
// 	maxEval  = +10000
// 	minEval  = -maxEval
// 	mateEval = maxEval + 1
// 	noScore  = minEval - 1
// )

// var pieceVal = [16]int{100, -100, 300, -300, 350, -350, 500, -500, 1000, -1000, 10000, -10000, 0, 0, 0, 0}

// var knightPosValue = [8]int{-4, -3, -2, +2, +2, 0, -2, -4}
// var knightRank = [8]int{-15, 0, +5, +6, +7, +8, +2, -4}
// var centerFile = [8]int{-8, -1, 0, +1, +1, 0, -1, -3}
// var kingFile = [8]int{+1, +2, 0, -2, -2, 0, +2, +1}
// var kingRank = [8]int{+1, 0, -2, -4, -6, -8, -10, -12}
// var pawnRank = [8]int{0, 0, 0, 0, +2, +6, +25, 0}
// var pawnFile = [8]int{0, 0, +1, +10, +10, +8, +10, +8}

// const longDiag = 10

// var knightPosValue = [64]int{}
// var knightRank = [64]int{}
// var centerFile = [64]int{}
// var kingFile = [64]int{}
// var kingRank = [64]int{}
// var pawnRank = [64]int{}
// var pawnFile = [64]int{}

// // Piece Square Table
var pSqTab2 [12][64]int

// Evaluate (new function)
func Evaluate(b *boardStruct) int {
	ev := 0
	// for sq := A1; sq <= H8; sq++ {
	// 	p12 := b.sq[sq]
	// 	if p12 == empty {
	// 		continue
	// 	}
	// 	ev += pieceVal[p12]
	// 	ev += pcSqScore(p12, sq)
	// }

	return ev
}

// // --------------------------------------------------- //

// // Score returns the pc2pt square table value for a given pc2pt on a given square. Stage = MG/EG
// func pcSqScore(p12, sq int) int {
// 	return pSqTab[p12][sq]
// }

// PstInit inits the pieces-square-tables when the program starts
func pcSqInit2() {
	tell("info string pStInit starter")
	pawnScore := [64]int{
		90, 90, 90, 90, 90, 90, 90, 90,
		30, 30, 30, 40, 40, 30, 30, 30,
		20, 20, 20, 30, 30, 30, 20, 20,
		10, 10, 10, 20, 20, 10, 10, 10,
		5, 5, 10, 20, 20, 5, 5, 5,
		0, 0, 0, 5, 5, 0, 0, 0,
		0, 0, 0, -10, -10, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	knightScore := [64]int{
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 10, 10, 0, 0, -5,
		-5, 5, 20, 20, 20, 20, 5, -5,
		-5, 10, 20, 30, 30, 20, 10, -5,
		-5, 10, 20, 30, 30, 20, 10, -5,
		-5, 5, 20, 10, 10, 20, 5, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, -10, 0, 0, 0, 0, -10, -5,
	}
	bishopScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 10, 10, 0, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 10, 0, 0, 0, 0, 10, 0,
		0, 30, 0, 0, 0, 0, 30, 0,
		0, 0, -10, 0, 0, -10, 0, 0,
	}
	rookScore := [64]int{
		50, 50, 50, 50, 50, 50, 50, 50,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 10, 20, 20, 10, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
	}
	kingScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 5, 5, 5, 0, 0,
		0, 5, 5, 10, 10, 5, 5, 0,
		0, 5, 10, 20, 20, 10, 5, 0,
		0, 5, 10, 20, 20, 10, 5, 0,
		0, 0, 5, 10, 10, 5, 0, 0,
		0, 5, 5, -5, -5, 0, 5, 0,
		0, 0, 5, 0, -15, 0, 10, 0,
	}
	queenScore := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	pSqTab[wP] = pawnScore
	pSqTab[wN] = knightScore
	pSqTab[wR] = rookScore
	pSqTab[wK] = kingScore
	pSqTab[wB] = bishopScore
	pSqTab[wQ] = queenScore

	// for Black
	for pc := Pawn; pc <= King; pc++ {

		wP12 := pt2pc(pc, WHITE)
		bP12 := pt2pc(pc, BLACK)

		for bSq := 0; bSq < 64; bSq++ {
			wSq := oppRank(bSq)
			pSqTab[bP12][bSq] = -pSqTab[wP12][wSq]
		}
	}
}

// mirror the rank_sq
func oppRank2(sq int) int {
	fl := sq % 8
	rk := sq / 8
	rk = 7 - rk
	return rk*8 + fl
}
