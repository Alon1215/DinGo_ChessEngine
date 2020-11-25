# DinGo (chess engine)

![Dingo](https://github.com/Alon1215/DinGo_chess.engine/blob/master/.idea/rsz_1rsz_1rsz_1rsz_687457.png?raw=true)


Chess engine (following UCI protocol) fully implemented in Go (Golang).
Built as a part of [Tomer Gonen](https://github.com/yodatk)'s and my chess engines competition which will take place at 30.10.20 . 

The engine will be featuring state of the art search techniques and based on bitboards.  

### board representation:
* bitboards
* magic numbers used for bitBoards move generation.

### search methods:
* alpha-beta search
* Parallel search
* Iterative Deepening
* Transposition Table
* Null Move Pruning
* Late Move Reductions

***Credits for the the comprehensive guidance to [Caro Kanns](https://www.youtube.com/playlist?list=PLftcy-r3mehgu4gikLTFoI1CXh2bHm3rf), [Chizhov Vadim](https://github.com/ChizhovVadim/CounterGo), and [Chess Programming Wiki](https://www.chessprogramming.org/Main_Page).***

Future goals:
* Principal Variation Search
* Improve evaluation
* Implementing deep learning techniques
* Upload to lichess.com

### How to install & play:
* Install a Chess GUI which follows the UCI protocol (for example: Arena, Scid).
* Use the command "go build" to compile the engine (DinGo.exe is already attached to the repository, in case Go is not installed on the computer).
* In Arena (or other GUI program) go to: Engines-> Manage Engines -> Install new engine. Install the DinGo.exe file.
* Dingo is installed, and is able to compete.
* Soon: live bot available (in lichess.com).
