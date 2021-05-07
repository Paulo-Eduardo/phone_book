package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func New() *sql.DB {
	DbConn, err := sql.Open("mysql", fmt.Sprintf("root:password123@tcp(%s)/phonebookdb", os.Getenv("DB_HOST")))
	if err != nil {
		log.Fatal(err)
	}
	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)

	return DbConn
}
