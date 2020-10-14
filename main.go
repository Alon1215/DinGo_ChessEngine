package main

func main() {
	tell("info string Hello DinGo")

	uci(input())

	tell("info string quits DinGo")
}

// ----------------------

func init() {
	initFenSq2Int()
	initMagic()
	initAtksKings()
	initAtksKnights()
	initCastlings()
	pSqInit()
	board.newGame()
}
