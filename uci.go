package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var tell = mainTell // set default tell
var trim = strings.TrimSpace
var low = strings.ToLower

func uci(frGUI chan string, myTell func(text ...string)) {
	tell = myTell
	tell("info string Hello from uci")
	frEng, toEng := engine()
	quit := false
	cmd := ""
	words := []string{}
	bm := ""
	for !quit {
		select {
		case cmd = <-frGUI:
			words = strings.Split(cmd, " ")
		case bm = <-frEng:
			handleBm(bm)
			continue
		}
		words[0] = trim(low(words[0]))
		switch words[0] {
		case "uci":
			handleUci()
		case "setoption":
			handleSetOption(words)
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleIsNewgame()
		case "position":
			handlePosition(words)
		case "debug":
			handleDebug(words)
		case "register":
			handleIsRegister(words)
		case "go":
			handleGo(words)
		case "ponderhit":
			handlePonderhit()
		case "stop":
			handleStop(toEng, &bInfinite)
		case "quit", "q":
			handleQuit(toEng)
			//quit = true
			continue
		}

	}
}

func handleUci() {
	tell("id name DinGo")
	tell("id author Alon1215")

	tell("option name Hash type spin default 32 min 1 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}

func handleQuit(toEng chan string) {
	toEng <- "stop"
}

func handleBm(bm string, bInfinite bool) {
	if bInfinite {
		saveBm = bm
		return
	}
	tell(bm)
}

func handleStop(toEng chan string, bInfinite bool) {
	if *bInfinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}

		toEng <- "stop"
		*bInfinite = true
	}
}

// Not impleneted uci commands:
func handleSetOption(words []string) {
	tell("info string set option ", strings.Join(option, " "))
	tell("info string not implemented")
}

func handleNewgame(option []string) {
	tell("info string ucinewgame not implemented")
}

func handlePosition(words []string) {
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "startpos":
			tell("info string position startpos not implemented")
		case "fen":
			tell("info string position fen not implemented")
		default:
			tell("info string position ", words[1], " not implemented")
		}
	}
}

func handleGo(words []string) {
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "searchmoves":
			tell("info string go searchmoves not implemented")
		case "ponder":
			tell("info string go ponder not implemented")
		case "wtime":
			tell("info string go wtime not implemented")
		case "btime":
			tell("info string go btime not implemented")
		case "winc":
			tell("info string go winc not implemented")
		case "binc":
			tell("info string go binc not implemented")
		case "movestogo":
			tell("info string go movestogo not implemented")
		case "depth":
			tell("info string go depth not implemented")
		case "nodes":
			tell("info string go nodes not implemented")
		case "movetime":
			tell("info string go movetime not implemented")
		case "mate":
			tell("info string go mate not implemented")
		case "infinite":
			tell("info string go infinite not implemented")
		default:
			tell("info string go ", words[1], " not implemented")
		}
	} else {
		tell("info string go not implemented")
	}
}

func handlePonderhit() {
	tell("info string ponderhit not implemented")
}

func handleDebug(option []string) {
	tell("info string debug not implemented")
}

func handleRegister(option []string) {
	tell("info string register not implemented")
}

func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	go func() {
		var reader *bufio.Reader
		reader = bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}
