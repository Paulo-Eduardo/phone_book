package healthcheck

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Paulo-Eduardo/phone_book/cors"
	"github.com/Paulo-Eduardo/phone_book/logger"
)

func SetupRoutes(apiBasePath string) {
	handleHealthCheck := http.HandlerFunc(HealthCheckHandler)
	http.Handle(fmt.Sprintf("%s/health-check", apiBasePath), logger.Middleware(cors.Middleware(handleHealthCheck)))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
