package main

import (
	"encoding/json"
	"log"
	"math/rand"
)

type Color uint

type Board struct {
	Diamonds bool `json:"diamonds"`
	NumRows uint `json:"num_rows"`
	NumCols uint `json:"num_cols"`
	Data     [][]Color `json:"data"`
}

const (
	TOP_PLAYER = 0
	BOT_PLAYER = 1
)

// Generates a board of diamonds
func CreateDiamondBoard(num_colors uint, num_rows uint, num_cols uint) Board {
	board := make([][]Color, num_rows)
	for i := 0; uint(i) < (num_rows-1)/2; i++ {
		// Non-offset rows
		board[2*i] = make([]Color, num_cols)
		for j := range num_cols {
			board[2*i][j] = Color(rand.Intn(int(num_colors)))
		}

		// Offset rows
		board[2*i+1] = make([]Color, num_cols-1)
		for j := range num_cols - 1 {
			board[2*i+1][j] = Color(rand.Intn(int(num_colors)))
		}
	}
	// Last row

	// Players must have different colors
	top_player_color := board[0][num_cols-1]
	bot_player_color := Color(rand.Intn(int(num_colors - 1)))
	if bot_player_color > top_player_color {
		board[num_rows-1][0] = bot_player_color + Color(1)
	} else {
		board[num_rows-1][0] = bot_player_color
	}

	for j := 1; j < int(num_cols); j++ {
		board[num_rows-1][j] = Color(rand.Intn(int(num_colors)))
	}

	return Board{true, num_rows, num_cols, board}
}

// Gets neighbors of a given position on a diamond board
func (board *Board) getNeighbors(pos [2]uint) [][2]uint {
	// Not exactly elegant, but it's good enough
	i := pos[0]
	j := pos[1]

	if i&1 != 0 {
		return ([][2]uint{
			{i - 1, j},
			{i - 1, j + 1},
			{i + 1, j},
			{i + 1, j + 1},
		})
	}

	result := make([][2]uint, 0, 4)
	if i > 0 && j > 0 {
		result = append(result, ([2]uint{i - 1, j - 1}))
	}
	if i > 0 && j < board.NumCols-1 {
		result = append(result, [2]uint{i - 1, j})
	}
	if i < board.NumRows-1 && j > 0 {
		result = append(result, [2]uint{i + 1, j - 1})
	}
	if i < board.NumRows-1 && j < board.NumCols-1 {
		result = append(result, [2]uint{i + 1, j})
	}
	return result
}

// Performs flood-fill on the board, changing the players color (i.e. making a move)
func (board *Board) ChangePlayerColor(player uint8, color Color) {
	var start_pos [2]uint
	if player == TOP_PLAYER {
		start_pos = [2]uint{0, board.NumCols - 1}
	} else {
		start_pos = [2]uint{board.NumRows - 1, 0}
	}

	last_color := board.Data[start_pos[0]][start_pos[1]]
	board.Data[start_pos[0]][start_pos[1]] = color

	to_visit := make([][2]uint, 1)
	to_visit[0] = start_pos
	for len(to_visit) != 0 {
		// Pop position
		pos := to_visit[len(to_visit)-1]
		to_visit = to_visit[:len(to_visit)-1]

		neighbors := board.getNeighbors(pos)
		for len(neighbors) != 0 {
			// Pop neighbor
			neighbor_pos := neighbors[len(neighbors)-1]
			neighbors = neighbors[:len(neighbors)-1]

			if board.Data[neighbor_pos[0]][neighbor_pos[1]] == last_color {
				board.Data[neighbor_pos[0]][neighbor_pos[1]] = color
				to_visit = append(to_visit, neighbor_pos)
			}
		}
	}
}

func Decode(boardEncoded string) Board {
	var board Board
	err := json.Unmarshal([]byte(boardEncoded), &board)
	if err != nil {
		log.Println("Board Unmarshal: ", err)
	}

	return board
}

func (board *Board) Encode() []byte {
	boardEncoded, err := json.Marshal(board)
	if err != nil {
		log.Println("Board Marshal: ", err)
		return nil
	}

	return boardEncoded
}
