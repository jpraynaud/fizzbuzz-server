package main

import (
	"encoding/json"
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

// pageResponse represents an error response
type pageResponse struct {
	Error    bool        `json:"error"`
	Response interface{} `json:"response"`
}

// Error page
func errorPage(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Errorf("%s - %s - %d - %s", r.Method, r.RequestURI, status, message)
	w.WriteHeader(status)
	pageResponse := pageResponse{true, message}
	json.NewEncoder(w).Encode(pageResponse)
}

// Handle FizzBuzz render
func renderHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare input parameters
	vars := r.URL.Query()
	limit, err := strconv.Atoi(vars.Get("limit"))
	if err != nil {
		errorPage(w, r, http.StatusBadRequest, "limit parameter must be an integer\n")
		return
	}
	int1, err := strconv.Atoi(vars.Get("int1"))
	if err != nil {
		errorPage(w, r, http.StatusBadRequest, "int1 parameter must be an integer\n")
		return
	}
	int2, err := strconv.Atoi(vars.Get("int2"))
	if err != nil {
		errorPage(w, r, http.StatusBadRequest, "int2 parameter must be an integer\n")
		return
	}
	str1 := vars.Get("str1")
	str2 := vars.Get("str2")

	// Render request
	renderer := render.NewRenderer()
	request := render.NewRequest(limit, int1, int2, str1, str2)
	response := renderer.Render(request)
	if err := response.Error; err != nil {
		errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Write response
	first := true
	for line := range response.Lines {
		if !first {
			line = fmt.Sprintf(",%s", line)
		}
		fmt.Fprint(w, line)
		first = false
	}
}
	// TODO
}

// Handles FizzBuzz rendering statistics
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
