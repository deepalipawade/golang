package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

//websocket - full duplex over single tcp

var conn *websocket.Conn

func main(){
  router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":9000")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "game.html")
	// homeTemplate.Execute(w, "ws://"+r.Host+"/")
}

func newGameHandler(rw http.ResponseWriter, req *http.Request) {
	// word := "elephant"
//
	start := time.Now()
	type Score struct {
		name string
		time []int
	}
	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	_, level, err := c.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}

	c.WriteMessage(websocket.TextMessage, []byte(getTopScores()))

	fmt.Println(level)
	var NewSudoku Sudoku
	NewSudoku.grid = make([][]int, 9)
	NewSudoku.createGrid()
	NewSudoku.RandomizeSudoku()
	NewSudoku.SudokuGenrator() // backtrack
	NewSudoku.LevelWiseSudoku(string(level))
	str := NewSudoku.getStringArray()
	c.WriteMessage(websocket.TextMessage, []byte(str))

	for {
		// score := Score{}
		var userData map[string]int

		_, recvData, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		//Extracting data from UI
		_ = json.Unmarshal(recvData, &userData)
		value := userData["value"]
		row := userData["row"]
		col := userData["col"]

		//To Do - create func for directly checking by comparing grid position value and user entered value
		blockCheck := violate(NewSudoku.userGrid, row, col, value)
		if blockCheck {
			c.WriteMessage(websocket.TextMessage, []byte("violation"))
		} else {
			NewSudoku.userGrid[row][col] = value
			win := NewSudoku.checkWin()
			if win {
				c.WriteMessage(websocket.TextMessage, []byte("win"))
				userTiming := time.Since(start)
				// Getting player name
				_, nameData, _ := c.ReadMessage()
				name := string(nameData)
				saveScore(userTiming, name)
				break
			}
		}
	}


}

func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))
	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)
	return
}

///// TEMP:


//grid.go
