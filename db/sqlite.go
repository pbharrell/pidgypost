package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var logger *log.Logger

func init() {
	f, err := os.OpenFile("database.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "", log.LstdFlags|log.Lshortfile)

	// Create or open the SQLite file
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create messages table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT CHECK( type IN ('Incoming','Outgoing') ) NOT NULL,
		sender TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		logger.Fatalf("Failed to create table: %v", err)
	}

	logger.Println("Database initialized and 'messages' table created.")
}
