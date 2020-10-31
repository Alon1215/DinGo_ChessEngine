package db

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// // GetRes request from DB for best pv
// type GetRes struct {
// 	Status string   `json:"status"`
// 	Score  int      `json:"score"`
// 	Depth  int      `json:"depth"`
// 	Pv     []string `json:"pv"`
// 	PvSAN  []string `json:"pvSAN"`
// }

// func get() ( ok bool) {
// 	fmt.Println("1. Performing Http Get...")
// 	// resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
// 	resp, err := http.Get(fmt.Sprintf("https://www.chessdb.cn/cdb.php?action=querypv&board=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&json=1") "https://jsonplaceholder.typicode.com/todos/1")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	defer resp.Body.Close()
// 	bodyBytes, _ := ioutil.ReadAll(resp.Body)

// 	// Convert response body to string
// 	bodyString := string(bodyBytes)
// 	fmt.Println("API Response as String:\n" + bodyString)

// 	// Convert response body to Todo struct
// 	var getStruct GetRes
// 	json.Unmarshal(bodyBytes, &getStruct)
// 	//fmt.Printf("API Response as struct %+v\n", todoStruct)
// }
