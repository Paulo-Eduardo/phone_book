package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Paulo-Eduardo/phone_book/database"
	"github.com/Paulo-Eduardo/phone_book/phonebook"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	argsWithoutProg := os.Args[1:]
	dbConn := database.New()
	timeout, err := strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		log.Fatal("Timeout must be a integer")
	}
	phonebook.SetupRoutes(apiBasePath, dbConn, timeout)
	log.Println("Server runnint at port: " + argsWithoutProg[0])
	log.Fatal(http.ListenAndServe(":"+argsWithoutProg[0], nil))
}
