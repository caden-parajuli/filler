package main

import (
	"flag"
	"log"

	"github.com/gorilla/websocket"
	"net/http"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// TODO add config file support
var dbAddress = flag.String("dbaddr", "file:database/games.db", "Database address")
var addr = flag.String("addr", "localhost:42069", "HTTP service address")

var db *sql.DB
var upgrader = websocket.Upgrader{
	Subprotocols: []string{"JSON-v1"},
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	var err error
	db, err = sql.Open("sqlite3", *dbAddress)
	if err != nil {
		log.Fatal(err)
	}
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
