package DinGo_ChessEngine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var tell = mainTell

func uci(frGUI chan string, myTell func(text ...string)) {
	tell = myTell
	tell("info string Hello from uci")
	frEng, toEng := engine()
	quit := false
	cmd := ""
	bm := ""
	for quit == false {
		select {
		case cmd = <-frGUI:
		case bm = <-frEng:
			handleBm(bm)
			continue
		}
		switch cmd {
		case "uci":
		case "stop":
			handleStop(toEng)
		case "quit", "q":
			quit = true
			continue
		}

	}
}

func handleBm(bm string) {
	tell(bm)
}

func handleStop(toEng chan string) {
	toEng <- "stop"
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
}

func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}
