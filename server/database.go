// Functions that will need to be moved

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Creates a SQLite database file
func CreateFileSQLite() {
	if (*dbAddress)[:5] != "file:" {
		log.Fatal("Attempted to create database file for non-SQLite database")
	}
	filename := (*dbAddress)[5:]

	// Create parent folder if necessary
	err := os.MkdirAll(filepath.Dir(filename), 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal("Could not create SQLite database directory: ", err)
	}

	// Create the database file if it does not already exist
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("Could not create SQLite database: ", err)
	}
	err = file.Close()
	if err != nil {
		log.Println("Could not close SQLite database file: ", err)
	}
}

// Drops tables from the database
func DropTables(db *sql.DB) {
	_, err := db.Exec(createTablesCommands)
	if err != nil {
		log.Fatal("Could not drop tables: ", err)
	}

	log.Println("Dropped database tables")
}

func CreateTablesSQLite(db *sql.DB) {
	_, err := db.Exec(createTablesCommands)
	if err != nil {
		log.Fatal("Could not create tables: ", err)
	}
}

// map Client IDs to Client structs. Needed to determine
// the connection corresponding to a players opponent
var clientsLock = sync.RWMutex{}
var clients = map[string]*Client{}

// Creates a new client and updates database accordingly
// if the id given is not already present
func NewPlayer(client *Client, id string) {
	// Check if player is already in database
	if id != "" {
		var game uint64
		row := db.QueryRow(`SELECT gameId FROM players WHERE id = ?;`, id)
		err := row.Scan(&game)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Println("QueryRow: ", err)
			}
		} else {
			client.id = id
			client.game = game

			clientsLock.Lock()
			clients[client.id] = client
			clientsLock.Unlock()

			return
		}
	}

	// Create new ID/client
	id, err := gonanoid.New()
	if err != nil {
		log.Fatal("nanoid new: ", err)
	}
	client.id = id

	clientsLock.Lock()
	clients[client.id] = client
	clientsLock.Unlock()

	// Insert into players table
	_, err = db.Exec(`INSERT INTO players VALUES(?, ?);`, client.id, client.game)
	if err != nil {
		log.Print("INSERT players: ", err)
	}

	return
}

// Reads board, players, and whose turn it is from database
func GetBoard(gameId uint64) (board *Board, player1Id, player2Id, turnPlayerId string) {
	if gameId == NO_GAME {
		log.Println("Tried to get board for null game")
		return nil, "", "", ""
	}
	var boardEncoded string
	row := db.QueryRow(`SELECT player1Id, player2Id, board, turn FROM games WHERE id = ?;`, gameId)
	err := row.Scan(&player1Id, &player2Id, &boardEncoded, &turnPlayerId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Print("GetBoard: ", err)
		}
		return nil, "", "", ""
	}

	boardDecoded := Decode(boardEncoded)

	return &boardDecoded, player1Id, player2Id, turnPlayerId
}

// Sets database record for board, updating turn
func SetBoard(gameId uint64, board *Board, turnPlayerId string) {
	if gameId == NO_GAME {
		log.Println("Tried to set board for null game")
		return
	}

	boardEncoded := board.Encode()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Transaction: ", err)
	}

	_, err = tx.Exec("UPDATE games SET board = ?, turn = ? WHERE id = ?;", boardEncoded, turnPlayerId, gameId)
	if err != nil {
		tx.Rollback()
		log.Print("UPDATE games: ", err)
	}
	tx.Commit()
}

func NewGame(params GameParamsReq) (gameId uint64, board *Board, err error) {
	if !params.IsDiamonds {
		return NO_GAME, nil, errors.ErrUnsupported
	}

	board = CreateDiamondBoard(params.NumColors, params.NumRows, params.NumCols)
	boardEncoded := board.Encode()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Transaction: ", err)
	}

	row := tx.QueryRow("INSERT INTO games (player1Id, player2Id, board, turn) VALUES(?, '', ?, ?) RETURNING id;", params.Id, boardEncoded, params.Id)
	err = row.Scan(&gameId)
	if err != nil {
		tx.Rollback()
		log.Print("INSERT INTO games: ", err)
		return NO_GAME, nil, err
	}
	tx.Exec("UPDATE players SET gameId = ? WHERE id = ?", gameId, params.Id)

	tx.Commit()

	return gameId, board, nil
}

func TryJoin(id string, gameId uint64) (success bool, board *Board, opponent string, myTurn bool) {
	if gameId == NO_GAME {
		log.Println("Tried to join null game")
		return false, nil, "", false
	}

	var boardEncoded string
	var player1Id string
	var player2Id string
	var turnPlayerId string

	tx, err := db.Begin()
	if err != nil {
		log.Println("Transaction: ", err)
	}

	row := tx.QueryRow(`SELECT player1Id, player2Id, board, turn FROM games WHERE id = ?;`, gameId)
	err = row.Scan(&player1Id, &player2Id, &boardEncoded, &turnPlayerId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Print("GetBoard: ", err)
		}
		return false, nil, "", false
	}

	var turn bool
	if turnPlayerId == "" {
		_, err = tx.Exec(`UPDATE games SET player2Id = ?, turn = ? WHERE id = ?`, id, id, gameId)
		turn = true
	} else {
		_, err = tx.Exec(`UPDATE games SET player2Id = ? WHERE id = ?`, id, gameId)
		turn = false
	}
	if err != nil {
		log.Print("TryJoin: ", err)
		err = tx.Rollback()
		if err != nil {
			log.Print("TryJoin Rollback: ", err)
		}
	}
	tx.Commit()

	boardDecoded := Decode(boardEncoded)

	return true, &boardDecoded, player1Id, turn
}
