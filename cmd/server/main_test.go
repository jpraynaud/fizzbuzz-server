package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jpraynaud/fizzbuzz-server/pkg/render"
	log "github.com/sirupsen/logrus"
)

// validateHandler is a helper that validates a handler given a request, a wanted http code and a wanted http body
func validateHandler(t *testing.T, handler http.HandlerFunc, request *http.Request, code int, body apiResponse) {
	// Create handler recorder
	recorder := httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(recorder, request)
	// Validate recorder code and body
	if recorder.Code != code {
		t.Errorf("handler returned status code %v, want %v", recorder.Code, code)
	}
	got := strings.Trim(recorder.Body.String(), "\n")
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Trim(string(bodyJSON), "\n")
	if got != want {
		t.Errorf("handler returned body %q want %q", got, want)
	}
}

func Test_renderHandler(t *testing.T) {
	// Prepare tests data
	type args struct {
		handler           http.HandlerFunc
		method            string
		path              string
		query             string
		codeWanted        int
		apiResponseWanted apiResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{"Render Bad Request", args{renderHandler, "GET", "/render", "", http.StatusBadRequest, apiResponse{true, "limit parameter must be an integer, value  was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=Z&int1=Z&int2=Z&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "limit parameter must be an integer, value Z was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=20&int1=Z&int2=Z&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "int1 parameter must be an integer, value Z was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=Z&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "int2 parameter must be an integer, value Z was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=0&int1=3&int2=5&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "limit parameter must be >= 1, value 0 was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=20&int1=0&int2=5&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "int1 parameter must be >= 1, value 0 was given"}}},
		{"Render Bad Request", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=0&str1=A&str2=B", http.StatusBadRequest, apiResponse{true, "int2 parameter must be >= 1, value 0 was given"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=A&str2=B", http.StatusOK, apiResponse{false, "1,2,A,4,B,A,7,8,A,B,11,A,13,14,AB,16,17,A,19,B"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=AA&str2=BBB", http.StatusOK, apiResponse{false, "1,2,AA,4,BBB,AA,7,8,AA,BBB,11,AA,13,14,AABBB,16,17,AA,19,BBB"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=30&int1=2&int2=7&str1=AAA&str2=BBB", http.StatusOK, apiResponse{false, "1,AAA,3,AAA,5,AAA,BBB,AAA,9,AAA,11,AAA,13,AAABBB,15,AAA,17,AAA,19,AAA,BBB,AAA,23,AAA,25,AAA,27,AAABBB,29,AAA"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=30&int1=3&int2=5&str1=喂&str2=世界", http.StatusOK, apiResponse{false, "1,2,喂,4,世界,喂,7,8,喂,世界,11,喂,13,14,喂世界,16,17,喂,19,世界,喂,22,23,喂,世界,26,喂,28,29,喂世界"}}},
	}
	// Reset statistics
	renderer.Statistics = render.NewStatistics()
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request, err := http.NewRequest(tt.args.method, tt.args.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.args.query != "" {
				mappedQuery, err := url.ParseQuery(tt.args.query)
				if err != nil {
					t.Fatal(err)
				}
				request.URL.RawQuery = mappedQuery.Encode()
			}
			// Validate handler
			validateHandler(t, tt.args.handler, request, tt.args.codeWanted, tt.args.apiResponseWanted)
		})
	}
}

func Test_statisticsHandler(t *testing.T) {
	// Prepare tests data
	type args struct {
		handler           http.HandlerFunc
		method            string
		path              string
		query             string
		codeWanted        int
		apiResponseWanted apiResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=A&str2=B", http.StatusOK, apiResponse{false, "1,2,A,4,B,A,7,8,A,B,11,A,13,14,AB,16,17,A,19,B"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=AA&str2=BBB", http.StatusOK, apiResponse{false, "1,2,AA,4,BBB,AA,7,8,AA,BBB,11,AA,13,14,AABBB,16,17,AA,19,BBB"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=A&str2=B", http.StatusOK, apiResponse{false, "1,2,A,4,B,A,7,8,A,B,11,A,13,14,AB,16,17,A,19,B"}}},
		{"Statistics OK", args{statisticsHandler, "GET", "/statistics", "", http.StatusOK, apiResponse{false, render.RequestStatistic{Request: render.Request{Limit: 20, Int1: 3, Int2: 5, Str1: "A", Str2: "B"}, Total: 2}}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=AA&str2=BBB", http.StatusOK, apiResponse{false, "1,2,AA,4,BBB,AA,7,8,AA,BBB,11,AA,13,14,AABBB,16,17,AA,19,BBB"}}},
		{"Render OK", args{renderHandler, "GET", "/render", "limit=20&int1=3&int2=5&str1=AA&str2=BBB", http.StatusOK, apiResponse{false, "1,2,AA,4,BBB,AA,7,8,AA,BBB,11,AA,13,14,AABBB,16,17,AA,19,BBB"}}},
		{"Statistics OK", args{statisticsHandler, "GET", "/statistics", "", http.StatusOK, apiResponse{false, render.RequestStatistic{Request: render.Request{Limit: 20, Int1: 3, Int2: 5, Str1: "AA", Str2: "BBB"}, Total: 3}}}},
	}
	// Reset statistics
	renderer.Statistics = render.NewStatistics()
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request, err := http.NewRequest(tt.args.method, tt.args.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.args.query != "" {
				mappedQuery, err := url.ParseQuery(tt.args.query)
				if err != nil {
					t.Fatal(err)
				}
				request.URL.RawQuery = mappedQuery.Encode()
			}
			// Validate handler
			validateHandler(t, tt.args.handler, request, tt.args.codeWanted, tt.args.apiResponseWanted)
		})
	}
}

func Test_loggingSetup(t *testing.T) {
	// Prepare tests data
	tests := []struct {
		name        string
		environment string
		want        log.Level
	}{
		{"Development", "development", log.DebugLevel},
		{"Production", "production", log.InfoLevel},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment = tt.environment
			if got := loggingSetup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggingSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createRouter(t *testing.T) {
	// Prepare tests data
	type args struct {
		method     string
		path       string
		codeWanted int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Render", args{"GET", "/render", http.StatusBadRequest}},
		{"Statistics", args{"GET", "/statistics", http.StatusOK}},
	}
	// Prepare test server
	router := createRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request, err := http.NewRequest(tt.args.method, fmt.Sprintf("%s%s", server.URL, tt.args.path), nil)
			if err != nil {
				t.Fatal(err)
			}
			// Checks that server responds correctly
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tt.args.codeWanted {
				t.Errorf("server returned status code %v, want %v", response.StatusCode, tt.args.codeWanted)
			}
		})
	}
}
