// package new_sudoku
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Sudoku struct {
	grid  [][]int
	userGrid [][]int
	level int
}

func (s *Sudoku)getStringArray() string {
	var str string
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			str = str + strconv.Itoa(s.userGrid[i][j])
		}
	}
	return str
}

func RandomInt(MAX int, MIN int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(MAX-MIN+1) + MIN
}

func (NewSudoku *Sudoku)createGrid() {
	for j:=0; j<9 ; j++ {
		for i:=0; i<9; i++ {
			NewSudoku.grid[j] = append(NewSudoku.grid[j],0)
		}
	}
}


func (NewSudoku *Sudoku) LevelWiseSudoku(level string) {
	display("level wise grid :", NewSudoku.grid)
	// NewSudoku.userGrid = NewSudoku.grid
	grid := make([][]int , 9)

		for i:=0;i<9;i++{
			for j:=0;j<9;j++{
				grid[i] = append(grid[i],j)
			}
		}
	for i:=0;i<9;i++{
		for j:=0;j<9;j++{
			grid[i][j] = NewSudoku.grid[i][j]
		}
	}
	// copy(grid, NewSudoku.grid)
	// display("level wise original grid :", NewSudoku.grid)
	// display("level wise copy grid :", grid)
	//
	iteration := 10000
	itr, blanks := 0, 0
	switch level {
	case "0":
		blanks = RandomInt(8, 4)
	case "1":
		blanks = RandomInt(40, 20)
	case "2":
		blanks = RandomInt(60, 60)
	}
	// fmt.Println(blanks)
	i := 0
	for i = 0; i < iteration; i++ {
		row := (RandomInt(7832, 23)*i + iteration*(RandomInt(78, 2))) % 9
		col := (RandomInt(92812, 187)*i + iteration*(RandomInt(92, 8))) % 9
		if grid[row][col] != 0 {
			grid[row][col] = 0
			itr++
		}
		if itr == blanks {
			break
		}
	}
	// display("copy zeroes grid : ", grid)
	// display("original no zeroees grid : ", NewSudoku.grid)

	NewSudoku.userGrid = grid
	fmt.Println("iertatir : ", i)
	fmt.Println(" no of blanks : ", blanks)
}

// random but valid initializations
func (s *Sudoku) RandomizeSudoku() {
	val := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			row := (RandomInt(8, 0)*RandomInt(899, 76)*i + 1) % 9
			col := (RandomInt(9, 1)*RandomInt(974898, 8765)*j + 1) % 9

			//logic for increase randomization
			if i == 0 && j == 0 { //00
				val = (RandomInt(9, 1) * RandomInt(899, 76)) % 10
				s.grid[i][j] = val
			}
			if i == 0 && j != 0 { //00
				val = (RandomInt(9, 1) * RandomInt(974898, 8765) * j) % 10
			}
			if i != 0 && j == 0 { //00
				val = (RandomInt(9, 1) * RandomInt(89934, 76) * i) % 10
			} else { //00
				val = (RandomInt(9, 1) * RandomInt(899, 786) * i * j) % 10
			}
			//only assign if valid else skip
			if !violate(s.grid, row, col, val) {
				s.grid[row][col] = val
				// fmt.Println(row, col , val)
			}
		}
	}
}

//backtacking function for sudoku generation
func (s *Sudoku) SudokuGenrator() bool {
	row, col, empty := find_empty(s.grid)
	if empty == -1 {
		return true
	} else {
		// start := RandomInt(9,1)
		for i := 1; i <= 9; i++ {
			if !violate(s.grid, row, col, i) {
				s.grid[row][col] = i
				if s.SudokuGenrator() {
					return true
				}
				s.grid[row][col] = 0
			}
		}
	}
	return false

}

//find next empty location if not returns -1
func find_empty(grid [][]int) (i, j, err int) {
	for i = 0; i < 9; i++ {
		for j = 0; j < 9; j++ {
			if grid[i][j] == 0 {
				return i, j, err
			}
		}
	}
	err, i, j = -1, -1, -1
	return
}

// checking violation given : row col and its value
func violate(grid [][]int, row int, col int, val int) bool {
	// row col check
	for i := 0; i < 9; i++ {
		if i != col && grid[row][i] == val {
			return true
		}
		if i != row && grid[i][col] == val {
			return true
		}
	}
	//box check
	x := (row / 3) * 3
	y := (col / 3) * 3
	for i := x; i < x+3; i++ {
		for j := y; j < y+3; j++ {
			if val == grid[i][j] {
				return true
			}
		}
	}
	return false
}

var MAX, MIN = 9, 1

func DiagonalFill(globalSudoku [][]int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			globalSudoku[i][j] = 0
		}
	}

	// assign diagonal grids value first ...so we only call subgrid check function
	// no nned to check row column
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if i == j {
				for r := i; r < i+3; r++ {
					for c := j; c < j+3; c++ {
						globalSudoku[r][c] = rand.Intn(MAX-MIN+1) + MIN
						if !IsSudokuValid(globalSudoku) {
							c--
						}
					} // for c
				} // for r
			}
		}
	}
}

func display(msg string, grid [][]int) {
	fmt.Println(msg, " ")
	for i, v := range grid {
		fmt.Println(i, " : ", v)
	}
	fmt.Println("\n\n")
}

// For checking if user won or lost
func (s *Sudoku) checkWin() bool {
	display("grid : ", s.grid)
	display("userGrid ", s.userGrid)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.grid[i][j] != s.userGrid[i][j] {
				return false
			}
		}
	}
	return true
}
// valid check of complete sudoku
func IsSudokuValid(grid [][]int) bool {
	if IsRowAndColumnUnique(grid) && IsSubGridValid(grid) {
		return true
	}
	return false
}

// check for 3*3 subgrids for duplicate entries
func IsSubGridValid(grid [][]int) (valid bool) {
	valid = true
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			repeat := [10]int{}
			for r := i; r < i+3; r++ {
				for c := j; c < j+3; c++ {

					val := grid[r][c]
					if val != 0 {
						if repeat[val] == 0 {
							repeat[val]++
						} else {
							valid = false
							return
						}
					}
				} // for c
			} // for r
		} // for j
	} // fori
	return
}

func IsRowAndColumnUnique(grid [][]int) (valid bool) {
	valid = true
	for i := 0; i < 9; i++ {
		// reinitialiize all frequency values to 0
		rowgrid := [10]int{}
		colgrid := [10]int{}
		for j := 0; j < 9; j++ {
			// row-wise check
			val := grid[i][j]
			if val != 0 {
				if rowgrid[val] == 0 {
					rowgrid[val]++
				} else {
					valid = false
					return
				} //if

			}
			// col-wise check
			val = grid[j][i]
			if val != 0 {
				if colgrid[val] == 0 {
					colgrid[val]++
				} else {
					valid = false
					return
				} //if
			}
		} //for j
	} //for i
	return
}



// import (
// )

// var DB_USER =

var upgrader = websocket.Upgrader{}

var DB_USER string = os.Getenv("DATABASE_USERNAME")
var DB_PASSWORD string = os.Getenv("DATABASE_PASSWORD")
var DB_NAME string = os.Getenv("DATABASE_NAME")

var difficultLevel = map[string]int{
	"0": 30,
	"1": 50,
	"2": 70,
}
// To save score in database
func saveScore(userTime time.Duration, name string) {
	hours := int(userTime / time.Hour)
	minutes := int(userTime / time.Minute)
	seconds := int(userTime / time.Second)
	seconds = seconds - minutes*60
	current := time.Now()
	date := current.Format("2006-01-02")
	usertime := strconv.Itoa(hours) + ":" + strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO Scores(Name, Time, Date) VALUES (?,?,?)"

	insert, err := db.Query(sql, name, usertime, date)

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

type Score struct {
	Name string `json:"Name"`
	Time string `json:"Time"`
}

func getTopScores() string {
	var top []Score
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT Name, Time FROM Scores Order by Time LIMIT 5")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var tag Score
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Name, &tag.Time)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		top = append(top, tag)
	}
	jsonData, err := json.Marshal(top)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonData)
}

// func main() {
	// 	var NewSudoku Sudoku
	// 	NewSudoku.grid = make([][]int, 9)
	// 	NewSudoku.level = 1
	//
	// 	// this sudoku is just for reference used in checking validation logic
	// NewSudoku.grid[0] = []int{1, 5, 8, 2, 7, 6, 4, 9, 3}
	// NewSudoku.grid[1] = []int{3, 6, 7, 1, 4, 9, 8, 5, 2}
	// NewSudoku.grid[2] = []int{2, 4, 9, 3, 8, 5, 6, 1, 7}
	// NewSudoku.grid[3] = []int{6, 1, 4, 7, 3, 8, 9, 2, 5}
	// NewSudoku.grid[4] = []int{5, 8, 3, 9, 2, 1, 7, 6, 4}
	// NewSudoku.grid[5] = []int{7, 9, 2, 6, 5, 4, 1, 3, 8}
	// NewSudoku.grid[6] = []int{8, 2, 1, 5, 6, 7, 3, 4, 9}
	// NewSudoku.grid[7] = []int{9, 7, 5, 4, 1, 3, 2, 8, 6}
	// NewSudoku.grid[8] = []int{4, 3, 6, 8, 9, 2, 5, 7, 1}
	//
	// 	// NewSudoku.grid[0] = []int{1, 5, 0, 2, 7, 0, 4, 9, 3}
	// 	// NewSudoku.grid[1] = []int{3, 0, 7, 1, 4, 0, 8, 5, 0}
	// 	// NewSudoku.grid[2] = []int{0, 4, 9, 3, 8, 5, 6, 1, 7}
	// 	// NewSudoku.grid[3] = []int{6, 1, 4, 7, 3, 8, 0, 2, 5}
	// 	// NewSudoku.grid[4] = []int{5, 8, 0, 9, 2, 1, 7, 6, 4}
	// 	// NewSudoku.grid[5] = []int{7, 9, 2, 0, 5, 0, 1, 3, 0}
	// 	// NewSudoku.grid[6] = []int{0, 2, 1, 5, 6, 7, 3, 4, 9}
	// 	// NewSudoku.grid[7] = []int{9, 7, 5, 0, 1, 3, 2, 0, 6}
	// 	// NewSudoku.grid[8] = []int{0, 3, 6, 0, 9, 2, 5, 7, 0}
	//
	// 	//complete grid initialiize to zero
	// 	for i := 0; i < 9; i++ {
		// 		for j := 0; j < 9; j++ {
			// 			NewSudoku.grid[i][j] = 0
			// 		}
			// 	}
			//
			// 	display("initialiize to zero ", NewSudoku.grid)
			// 	// DiagonalFill(grid)
			// 	// display("diagonal sudoku  : ", NewSudoku.grid)
			//
			// 	RandomizeSudoku(NewSudoku.grid)
			// 	display("Random start sudoku  : ", NewSudoku.grid)
			//
			// 	SudokuGenrator(NewSudoku.grid)
			// 	fmt.Println("Is it valid ? : ", IsSudokuValid(NewSudoku.grid))
			// 	display("complete sudoku  : ", NewSudoku.grid)
			//
			// 	// according to level fill number of zeroes
			// 	level := 4
			// 	LevelWiseSudoku(NewSudoku.grid, level)
			// 	display("\nL4  sudoku  : ", NewSudoku.grid)
			// 	fmt.Println("Is it valid ? : ", IsSudokuValid(NewSudoku.grid))
			// }
