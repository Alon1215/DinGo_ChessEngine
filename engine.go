package DinGo_ChessEngine

func engine() (frEng, toEng chan string){
	tell("string info Hello from engine")
	go func() {
		for cmd := range toEng{
			switch cmd {
			case "stop":
			case "quit":

			}
		}
	}()
	return
}