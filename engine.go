package main

import (
	"fmt"
	"math"
	"time"
)

type searchLimits struct {
	depth     int
	nodes     uint64
	moveTime  int // in milliseconds
	infinite  bool
	startTime time.Time
	lastTime  time.Time

	//////////////// current //////////
	stop bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999
	s.infinite = false
	s.stop = false
}

func (s *searchLimits) setStop(st bool) {
	s.stop = st

}
func (s *searchLimits) setDepth(d int) {
	s.depth = d
}

func (s *searchLimits) setMoveTime(m int) {
	s.moveTime = m
}

func (s *searchLimits) setInfinite(b bool) {
	s.infinite = b
}

func engine() (toEngine chan bool, frEngine chan string) {
	tell("info string stopHello from engine")
	frEngine = make(chan string)
	toEngine = make(chan bool)
	// go func() {
	// 	for cmd := range toEng {
	// 		switch cmd {
	// 		case "stop":
	// 		case "quit":
	// 		case "go":
	// 			tell("info string Im thinking")
	// 			// TODO start the thinking process in the engine from "go"
	// 		}
	// 	}
	// }()

	go root(toEngine, frEngine)
	return
}

func root(toEngine chan bool, frEngine chan string) {
	b := &board
	ml := moveList{}
	for _ = range toEngine {
		tell("info string engine got go! X")
		ml = moveList{}
		genAndSort(b, &ml)
		for _, mv := range ml {
			b.move(mv)
			score := -search(b)
			b.unmove(mv)

			mv.packMove(adjEval(b, score))
		}
		ml.sort()
		tell("info score cp ", fmt.Sprintf("%v", ml[0].eval()), " depth 1 pv ", ml[0].String())
		frEngine <- fmt.Sprintf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
	}

}

func search(b *boardStruct) int {
	//return evaluate()
	return -1
}

func genAndSort(b *boardStruct, ml *moveList) {
	b.genAllMoves(ml)
	//TOD: update genAllMoves()
	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		v = adjEval(b, v)
		(*ml)[ix].packEval(v)
	}
	ml.sort()
}

func adjEval(b *boardStruct, ev int) int {
	if b.stm == BLACK {
		return -ev
	}
	return ev
}
