package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"net/http"

	_ "embed"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed database/create-tables.sql
var createTablesCommands string

// TODO add config file support

var dbEngine = flag.String("dbengine", "sqlite3", "Database engine")
var dbAddress = flag.String("dbaddr", "file:database/games.db", "Database address")

// Conveniently this port is not reserved
var addr = flag.String("addr", "localhost:42069", "HTTP service address")

var db *sql.DB

func main() {
	flag.Parse()
	log.SetFlags(0)

	openDB()
	defer db.Close()

	// Ping to ensure we connected
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected to database")

	// Start websocket server
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func openDB() {
	if *dbEngine == "sqlite3" {
		createDBSQLite()

		var err error
		db, err = sql.Open(*dbEngine, *dbAddress)
		if err != nil {
			log.Fatal(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("openDB sqlite3: ", err)
		}
		defer tx.Rollback()

		_, err = tx.Exec(createTablesCommands)
		if err != nil {
			log.Fatal("Could not create tables: ", err)
		}

		if err = tx.Commit(); err != nil {
			log.Fatal("Could not commit query to create tables: ", err)
		}

		log.Println("Finished with", createTablesCommands)

		return
	}

	var err error
	db, err = sql.Open(*dbEngine, *dbAddress)
	if err != nil {
		log.Fatal(err)
	}
}

// Creates a SQLite database file
func createDBSQLite() {
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
