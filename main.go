package main

import "strings"

const (
	name   = "DinGo"
	author = "Alon Michaeli"
)

func main() {
	tell("info string Hello DinGo")

	uci(input())

	tell("info string quits DinGo")
}

// ----------------------

func init() {
	initFen2Sq()
	initMagic()
	initKeys()
	initAtksKings()
	initAtksKnights()
	initCastlings()
	pcSqInit()
	board.newGame()
	handleSetOption(strings.Split("setoption name hash value 32", " "))
}
