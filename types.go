package main

import (
	"math"
	"time"
)

// SearchLimits sets limits for engine search
type SearchLimits struct {
	startTime time.Time
	lastTime  time.Time

	Ponder         bool
	Infinite       bool
	WhiteTime      int
	BlackTime      int
	WhiteIncrement int
	BlackIncrement int
	MoveTime       int
	MovesToGo      int
	Depth          int
	Nodes          int
	Mate           int
	Stop           bool
}

func (s *SearchLimits) init() {
	s.Depth = 9999
	s.Nodes = math.MaxInt64
	s.MoveTime = 3000
	s.Infinite = false
	s.Stop = false
	s.WhiteTime = -1
	s.BlackTime = -1
}
