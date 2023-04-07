package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	board       [][]string
	ships       []Ship
	shots       [][]bool
	numShots    int
	numAttempts int
}

type Ship struct {
	rowStart int
	colStart int
	length   int
	orient   string
	hits     int
}

func NewGame(boardSize int, shipSizes []int) *Game {
	// Initialize board
	board := make([][]string, boardSize)
	for i := 0; i < boardSize; i++ {
		board[i] = make([]string, boardSize)
		for j := 0; j < boardSize; j++ {
			board[i][j] = "-"
		}
	}

	// Place ships randomly
	ships := make([]Ship, len(shipSizes))
	for i, size := range shipSizes {
		var row, col int
		var orient string
		shipLength := size
		shipPlaced := false
		for !shipPlaced {
			row = randInt(boardSize)
			col = randInt(boardSize)
			orient = randOrient()
			if isValidShipPlacement(row, col, shipLength, orient, board) {
				placeShip(row, col, shipLength, orient, &board)
				ships[i] = Ship{row, col, shipLength, orient, 0}
				shipPlaced = true
			}
		}
	}

	// Initialize shots
	shots := make([][]bool, boardSize)
	for i := 0; i < boardSize; i++ {
		shots[i] = make([]bool, boardSize)
	}

	return &Game{board, ships, shots, 0, 0}
}

func (g *Game) ShowBoard() {
	fmt.Printf("   ")
	for i := 0; i < len(g.board); i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Printf("\n")
	for i, row := range g.board {
		fmt.Printf("%d  ", i)
		for _, cell := range row {
			switch cell {
			case "-":
				fmt.Printf("_ ")
			case "x":
				fmt.Printf("x ")
			case "+":
				fmt.Printf("+ ")
			case "S":
				fmt.Printf("S ")
			}
		}
		fmt.Printf("\n")
	}
}

func (g *Game) Attack(row int, col int) bool {
	if row < 0 || row >= len(g.board) || col < 0 || col >= len(g.board) {
		fmt.Println("Invalid attack coordinates")
		return false
	}

	if g.shots[row][col] {
		fmt.Println("You've already attacked this location")
		return false
	}

	g.shots[row][col] = true
	g.numAttempts++

	if g.board[row][col] == "-" {
		g.board[row][col] = "x"
		fmt.Println("Miss!")
		return false
	} else {
		for i, ship := range g.ships {
			if ship.orient == "H" {
				if row == ship.rowStart && col >= ship.colStart && col < ship.colStart+ship.length {
					g.ships[i].hits++
					g.numShots++
					if g.ships[i].hits == g.ships[i].length {
						for j := ship.colStart; j < ship.colStart+ship.length; j++ {
							g.board[row][j] = "+"
						}
						fmt.Println("Hit! You sunk my ship!")
						return true
					} else {
						g.board[row][col] = "+"
						fmt.Println("Hit!")
						return true
					}
				}
			} else {
				if col == ship.colStart && row >= ship.rowStart && row < ship.rowStart+ship.length {
					g.ships[i].hits++
					g.numShots++
					if g.ships[i].hits == g.ships[i].length {
						for j := ship.rowStart; j < ship.rowStart+ship.length; j++ {
							g.board[j][col] = "+"
						}
						fmt.Println("Hit! You sunk my ship!")
						return true
					} else {
						g.board[row][col] = "+"
						fmt.Println("Hit!")
						return true
					}
				}
			}
		}
	}
	return false
}

func (g *Game) IsOver() bool {
	for _, ship := range g.ships {
		if ship.hits < ship.length {
			return false
		}
	}
	return true
}

func main() {
	boardSize := 10
	shipSizes := []int{5, 4, 3, 3, 2}
	game := NewGame(boardSize, shipSizes)

	for !game.IsOver() {
		game.ShowBoard()
		fmt.Printf("Number of attempts: %d\n", game.numAttempts)
		fmt.Printf("Number of shots fired: %d\n", game.numShots)

		var row, col int
		fmt.Print("Enter row to attack: ")
		fmt.Scanln(&row)
		fmt.Print("Enter column to attack: ")
		fmt.Scanln(&col)

		game.Attack(row, col)
	}

	fmt.Println("Congratulations! You won!")
}

func isValidShipPlacement(row int, col int, length int, orient string, board [][]string) bool {
	if orient == "H" {
		if col+length > len(board) {
			return false
		}
		for i := col; i < col+length; i++ {
			if board[row][i] != "-" {
				return false
			}
		}
	} else {
		if row+length > len(board) {
			return false
		}
		for i := row; i < row+length; i++ {
			if board[i][col] != "-" {
				return false
			}
		}
	}
	return true
}

func placeShip(row int, col int, length int, orient string, board *[][]string) {
	if orient == "H" {
		for i := col; i < col+length; i++ {
			(*board)[row][i] = "S"
		}
	} else {
		for i := row; i < row+length; i++ {
			(*board)[i][col] = "S"
		}
	}
}

func randInt(n int) int {
	r := rand.Intn(n)
	return r
}

func randOrient() string {
	r := rand.Intn(2)
	if r == 0 {
		return "H"
	} else {
		return "V"
	}
}
