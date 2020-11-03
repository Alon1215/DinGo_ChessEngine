# DinGo (chess engine)

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

![Dingo](https://github.com/Alon1215/DinGo_chess.engine/blob/master/.idea/dingo-silhouette-fd8c25-md.png?raw=true)

