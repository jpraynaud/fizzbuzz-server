package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jpraynaud/fizzbuzz-server/pkg/render"
	log "github.com/sirupsen/logrus"
)

var renderer *render.Renderer

func init() {
	// Init FizzBuzz renderer
	renderer = render.NewRenderer()
}

var (
	environment string
	addr        string
)

func main() {
	// Parse flags
	flag.StringVar(&addr, "address", "0.0.0.0:8080", "server listening address")
	flag.StringVar(&environment, "environment", "development", "server environment (development or production)")
	flag.Parse()

	// Logging setup
	loggingSetup()

	// Start HTTP server
	router := createRouter()
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	log.Fatal(server.ListenAndServe())
}

// createRouter creates the router of the HTTP server
func createRouter() *mux.Router {
	log.WithFields(log.Fields{
		"environment": environment,
		"address":     addr,
	}).Info("Create server")
	router := mux.NewRouter()
	router.HandleFunc("/render", renderHandler).Methods(http.MethodGet)
	router.HandleFunc("/statistics", statisticsHandler).Methods(http.MethodGet)
	router.Use(loggingMiddleware)
	return router
}

// loggingSetup sets up logging
func loggingSetup() log.Level {
	log.SetOutput(os.Stdout)
	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
	return log.GetLevel()
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s - %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// apiResponse represents a page response
type apiResponse struct {
	Error    bool        `json:"error"`
	Response interface{} `json:"response"`
}

// Api error
func apiError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Errorf("%s - %s - %d - %s", r.Method, r.RequestURI, status, message)
	w.WriteHeader(status)
	apiResponse := apiResponse{true, message}
	json.NewEncoder(w).Encode(apiResponse)
}

// Handle FizzBuzz render
func renderHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare input parameters
	vars := r.URL.Query()
	limit, err := strconv.Atoi(vars.Get("limit"))
	if err != nil {
		apiError(w, r, http.StatusBadRequest, fmt.Sprintf("limit parameter must be an integer, value %s was given", vars.Get("limit")))
		return
	}
	int1, err := strconv.Atoi(vars.Get("int1"))
	if err != nil {
		apiError(w, r, http.StatusBadRequest, fmt.Sprintf("int1 parameter must be an integer, value %s was given", vars.Get("int1")))
		return
	}
	int2, err := strconv.Atoi(vars.Get("int2"))
	if err != nil {
		apiError(w, r, http.StatusBadRequest, fmt.Sprintf("int2 parameter must be an integer, value %s was given", vars.Get("int2")))
		return
	}
	str1 := vars.Get("str1")
	str2 := vars.Get("str2")

	// Render request
	request := render.NewRequest(limit, int1, int2, str1, str2)
	response := renderer.Render(request)
	if err := response.Error; err != nil {
		apiError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	lines := make([]string, 0)
	for line := range response.Lines {
		lines = append(lines, line)
	}

	// Write response
	apiResponse := apiResponse{false, strings.Join(lines, ",")}
	json.NewEncoder(w).Encode(apiResponse)
}

// Handles rendering statistics
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	statistics := renderer.Statistics
	topStatistic := statistics.TopStatistic()
	apiResponse := apiResponse{false, topStatistic}
	json.NewEncoder(w).Encode(apiResponse)
}
