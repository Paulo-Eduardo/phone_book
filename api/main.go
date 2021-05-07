package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Paulo-Eduardo/phone_book/database"
	"github.com/Paulo-Eduardo/phone_book/healthcheck"
	"github.com/Paulo-Eduardo/phone_book/phonebook"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	recordMetrics()

	argsWithoutProg := os.Args[1:]
	dbConn := database.New()
	timeout, err := strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		log.Fatal("Timeout must be a integer")
	}

	healthcheck.SetupRoutes(apiBasePath)
	phonebook.SetupRoutes(apiBasePath, dbConn, timeout)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server runnint at port: " + argsWithoutProg[0])
	log.Fatal(http.ListenAndServe(":"+argsWithoutProg[0], nil))
}
