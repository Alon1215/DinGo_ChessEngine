package main

func main() {
	tell("info string Hello DinGo")

	uci(input())

	tell("info string quits DinGo")
}

// ----------------------

func init() {
	initFen2Sq()
	initMagic()
	initAtksKings()
	initAtksKnights()
	initCastlings()
	pSqInit()
	board.newGame()
}
