package database

import (
	"database/sql"
	"log"
	"time"
)

func New() *sql.DB {
	DbConn, err := sql.Open("mysql", "root:password123@tcp(db)/phonebookdb")
	if err != nil {
		log.Fatal(err)
	}
	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)

	return DbConn
}
