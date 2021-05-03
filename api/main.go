package main

import (
	"log"
	"net/http"

	"github.com/Paulo-Eduardo/phone_book/database"
	"github.com/Paulo-Eduardo/phone_book/phonebook"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	phonebook.SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}