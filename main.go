package main

import (
	"strings"
)

const (
	name    = "DinGo"
	author  = "Alon Michaeli"
	version = "1.2"
)

func main() {
	tell("info string Hello DinGo")
	/////////////////

	// var protocol = &Protocol{
	// 	Name:    name,
	// 	Author:  author,
	// 	Version: version,
	// 	Engine:  engine,
	// 	Options: []uci.Option{
	// 		&uci.IntOption{Name: "Hash", Min: 4, Max: 1 << 16, Value: &engine.Hash},
	// 		&uci.IntOption{Name: "Threads", Min: 1, Max: runtime.NumCPU(), Value: &engine.Threads},
	// 		&uci.BoolOption{Name: "ExperimentSettings", Value: &engine.ExperimentSettings},
	// 	},
	// }

	// protocol.Run()

	/////////////////
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
	limits.init()
}
