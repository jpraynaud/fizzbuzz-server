// Package render handles requests that will be rendered according to the FizzBuzz algorithm (see README for details)
package render

import (
	"fmt"
	"sync"
)

// Request represents a request that will be rendered according to the FizzBuzz algorithm (see README for details)
// Limit is the number of lines that will be rendered (starting from 1 to Limit)
// Int1 (or Int2) represents the multiple of the line numbers that will display Str1 (or Str2) instead of their respective line number
type Request struct {
	Limit int    `json:"limit"`
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// NewRequest is the Request factory
func NewRequest(limit, int1, int2 int, str1, str2 string) *Request {
	return &Request{
		Limit: limit,
		Int1:  int1,
		Int2:  int2,
		Str1:  str1,
		Str2:  str2,
	}
}

// Validate checks that the request is valid and can be rendered by the FizzBuzz algorithm (see README for details)
// i.e. Limit/Int1/Int2 must be >= 1
func (r *Request) Validate() error {
	var err error = nil
	switch {
	case r.Limit < 1:
		err = fmt.Errorf("limit must be >= 1, value %d was given", r.Limit)
	case r.Int1 < 1:
		err = fmt.Errorf("int1 must be >= 1, value %d was given", r.Int1)
	case r.Int2 < 1:
		err = fmt.Errorf("int2 must be >= 1, value %d was given", r.Int2)
	}
	return err
}

// Response represents a response that will be returned when a request is rendered
type Response struct {
	Lines chan string
	Error error
}

// NewResponse is the Response factory
func NewResponse() *Response {
	return &Response{
		Lines: make(chan string),
		Error: nil,
	}
}

// Renderer represents a renderer for the FizzBuzz Algorithm (see README for details)
type Renderer struct {
	*Statistics
}

// NewRenderer is the Renderer factory
func NewRenderer() *Renderer {
	return &Renderer{
		Statistics: NewStatistics(),
	}
}

// Render renders the response associated with the request according to the FizzBuzz algorithm (see README for details)
func (rr *Renderer) Render(request *Request) *Response {
	response := NewResponse()
	if err := request.Validate(); err != nil {
		close(response.Lines)
		response.Error = err
		return response
	}
	go func() {
		for i := 1; i <= request.Limit; i++ {
			multiple := false
			line := ""
			if i%request.Int1 == 0 {
				multiple = true
				line += request.Str1
			}
			if i%request.Int2 == 0 {
				multiple = true
				line += request.Str2
			}
			if !multiple {
				line = fmt.Sprintf("%d", i)
			}
			response.Lines <- fmt.Sprint(line)
		}
		close(response.Lines)
	}()
	rr.Statistics.Record(request)
	return response
}

// Statistics represents statistics of requests rendering
type Statistics struct {
	Totals     map[Request]int
	TopRequest Request
	mutex      *sync.RWMutex
}

var statistics Statistics

// NewStatistics is the Statistics factory
func NewStatistics() *Statistics {
	if statistics.mutex == nil {
		statistics.Totals = make(map[Request]int)
		statistics.mutex = &sync.RWMutex{}
	}
	return &statistics
}

// Reset resets rendering statistics
func (s *Statistics) Reset() {
	s.mutex.Lock()
	s.Totals = make(map[Request]int)
	s.mutex.Unlock()
}

// Record records rendering statistics
func (s *Statistics) Record(request *Request) {
	s.mutex.Lock()
	s.Totals[*request]++
	if s.Totals[*request] > s.Totals[s.TopRequest] {
		s.TopRequest = *request
	}
	s.mutex.Unlock()
}

