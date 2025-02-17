package main

import (
	"flag"
	"log"

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
var dropDb = flag.Bool("drop-db", false, "Useful for testing. WARNING: this will drop the entire database before running!")

// Conveniently port 42069 is not reserved
var addr = flag.String("addr", "localhost:42069", "HTTP service address")


var db *sql.DB

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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
		CreateFileSQLite()

		var err error
		db, err = sql.Open(*dbEngine, *dbAddress)
		if err != nil {
			log.Fatal(err)
		}

		if *dropDb {
			DropTables(db)
		}
		CreateTablesSQLite(db)
		log.Println("Finished creating tables")

		return
	}

	var err error
	db, err = sql.Open(*dbEngine, *dbAddress)
	if err != nil {
		log.Fatal(err)
	}
}
