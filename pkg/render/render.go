// Package render handles requests that will be rendered according to the FizzBuzz algorithm (see README for details)
package render

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
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
	var err error
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

// Renderer represents the interface to render the FizzBuzz Algorithm (see README for details)
type Renderer interface {
	Render(ctx context.Context, request *Request) *Response
	StatisticRecorder
}

// Default renderer implementation
type renderer struct {
	*Statistics
}

// NewRenderer is the Renderer factory
func NewRenderer() Renderer {
	return &renderer{
		Statistics: NewStatistics(),
	}
}

// Render renders the response associated with the request according to the FizzBuzz algorithm (see README for details)
func (rr *renderer) Render(ctx context.Context, request *Request) *Response {
	defer rr.RecordStatistic(request)
	response := NewResponse()
	if err := request.Validate(); err != nil {
		defer close(response.Items)
		response.Error = err
		return response
	}
	log.Debugf("Request rendering started %+v", request)
	go func() {
		defer func() {
			log.Debugf("Request rendering done %+v", request)
			close(response.Items)
		}()
		for i := 1; i <= request.Limit; i++ {
			var item string
			switch {
			case i%request.Int1 == 0 && i%request.Int2 == 0:
				item = request.Str1 + request.Str2
			case i%request.Int1 == 0:
				item = request.Str1
			case i%request.Int2 == 0:
				item = request.Str2
			default:
				item = strconv.Itoa(i)
			}
			select {
			case response.Items <- item:
			case <-ctx.Done():
				log.Debugf("Request rendering cancelled %+v", request)
				return
			}
		}
	}()
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

// StatisticRecorder represents the interface for statistics recording of requests rendering
type StatisticRecorder interface {
	RecordStatistic(request *Request)
	GetStatistic(request *Request) *RequestStatistic
	GetTopStatistic() *RequestStatistic
	ResetStatistics()
}

// Statistics represents statistics of requests rendering
type Statistics struct {
	Totals     sync.Map
	TopRequest Request
}

// NewStatistics is the Statistics factory
func NewStatistics() *Statistics {
	return &Statistics{
		Totals: sync.Map{},
	}
}

// RecordStatistic records rendering statistics
func (s *Statistics) RecordStatistic(request *Request) {
	total, _ := s.Totals.Load(*request)
	totalI, _ := total.(int)
	totalI++
	totalTop, _ := s.Totals.Load(s.TopRequest)
	totalTopI, _ := totalTop.(int)
	s.Totals.Store(*request, totalI)
	if totalI > totalTopI {
		s.TopRequest = *request
	}
}

// GetStatistic returns rendering statistics of a request
func (s *Statistics) GetStatistic(request *Request) *RequestStatistic {
	total, _ := s.Totals.Load(*request)
	totalI, _ := total.(int)
	if totalI == 0 {
		return nil
	}
	return NewRequestStatistic(request, totalI)
}

// GetTopStatistic returns rendering statistics of the top request
func (s *Statistics) GetTopStatistic() *RequestStatistic {
	topRequest := s.TopRequest
	return s.GetStatistic(&topRequest)
}

// ResetStatistics resets all statistics currently recorded
func (s *Statistics) ResetStatistics() {
	s.Totals = sync.Map{}
	s.TopRequest = Request{}
}
