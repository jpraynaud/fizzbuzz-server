// Package render handles requests that will be rendered according to the FizzBuzz algorithm (see README for details)
package render

import (
	"fmt"
	"sync"
)

// Request represents a request that will be rendered according to the FizzBuzz algorithm (see README for details)
// Limit is the number of items that will be rendered (starting from 1 to Limit)
// Int1 (or Int2) represents the multiple of the item numbers that will display Str1 (or Str2) instead of their respective item number
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
		err = fmt.Errorf("limit parameter must be >= 1, value %d was given", r.Limit)
	case r.Int1 < 1:
		err = fmt.Errorf("int1 parameter must be >= 1, value %d was given", r.Int1)
	case r.Int2 < 1:
		err = fmt.Errorf("int2 parameter must be >= 1, value %d was given", r.Int2)
	}
	return err
}

// Response represents a response that will be returned when a request is rendered
type Response struct {
	Items chan string
	Error error
}

// NewResponse is the Response factory
func NewResponse() *Response {
	return &Response{
		Items: make(chan string),
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

		defer close(response.Items)
		response.Error = err
		return response
	}
	go func() {
		defer close(response.Items)
		for i := 1; i <= request.Limit; i++ {
			multiple := false
			item := ""
			if i%request.Int1 == 0 {
				multiple = true
				item += request.Str1
			}
			if i%request.Int2 == 0 {
				multiple = true
				item += request.Str2
			}
			if !multiple {
				item = fmt.Sprintf("%d", i)
			}
			response.Items <- item
		}
	}()
	rr.Statistics.Record(request)
	return response
}

// RequestStatistic represents the rendering statistics of a request
type RequestStatistic struct {
	Request `json:"request"`
	Total   int `json:"total"`
}

// NewRequestStatistic is the RequestStatistic factory
func NewRequestStatistic(request *Request, total int) *RequestStatistic {
	return &RequestStatistic{
		Request: *request,
		Total:   total,
	}
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
	return &Statistics{
		Totals: make(map[Request]int),
		mutex:  &sync.RWMutex{},
	}
}

// Record records rendering statistics
func (s *Statistics) Record(request *Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Totals[*request]++
	if s.Totals[*request] > s.Totals[s.TopRequest] {
		s.TopRequest = *request
	}
}

// Statistic returns rendering statistics of a request
func (s *Statistics) Statistic(request *Request) *RequestStatistic {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.Totals[*request] == 0 {
		return nil
	}
	return NewRequestStatistic(request, s.Totals[*request])
}

// TopStatistic returns rendering statistics of the top request
func (s *Statistics) TopStatistic() *RequestStatistic {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	topRequest := s.TopRequest
	if s.Totals[topRequest] == 0 {
		return nil
	}
	return s.Statistic(&topRequest)
}
