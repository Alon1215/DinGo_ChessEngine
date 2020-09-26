package main

func engine() (frEng, toEng chan string) {
	tell("info string stopHello from engine")
	frEng = make(chan string)
	toEng = make(chan string)
	go func() {
		for cmd := range toEng {
			switch cmd {
			case "stop":
			case "quit":

			}
		}
	}()
	return
}
