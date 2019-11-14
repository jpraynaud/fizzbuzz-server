package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jpraynaud/fizzbuzz-server/pkg/render"
	log "github.com/sirupsen/logrus"
)

var (
	environment string
	addr        string
)

func main() {
	// Parse flags
	flag.StringVar(&addr, "address", "0.0.0.0:8080", "server listening address")
	flag.StringVar(&environment, "environment", "development", "server environment (development or production)")
	flag.Parse()

	// Log setup
	log.SetOutput(os.Stdout)
	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	// Start HTTP server
	log.WithFields(log.Fields{
		"environment": environment,
		"address":     addr,
	}).Info("Start server")
	router := mux.NewRouter()
	router.HandleFunc("/render", renderHandler).Methods(http.MethodGet)
	router.HandleFunc("/statistics", statisticsHandler).Methods(http.MethodGet)
	router.Use(loggingMiddleware)
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	log.Fatal(server.ListenAndServe())
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s - %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Handle FizzBuzz render
func renderHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Handles FizzBuzz rendering statistics
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
