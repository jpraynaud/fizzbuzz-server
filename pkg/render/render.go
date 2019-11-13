// Package render handles requests that will be rendered according to the FizzBuzz algorithm (see README for details)
package render

import (
	"fmt"
)

// Request represents a request that will be rendered according to the FizzBuzz algorithm (see README for details)
// Limit is the number of lines that will be rendered (starting from 1 to Limit)
// Int1 (or Int2) represents the multiple of the line numbers that will display Str1 (or Str2) instead of their respective line number
type Request struct {
	Limit, Int1, Int2 int
	Str1, Str2        string
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
// i.e. Limit must be >= 0 and Int1/Int2 must be >= 1
func (r *Request) Validate() error {
	var err error = nil
	switch {
	case r.Limit < 0:
		err = fmt.Errorf("Limit must be >= 0, value %d was given", r.Limit)
	case r.Int1 < 1:
		err = fmt.Errorf("Int1 must be >= 1, value %d was given", r.Int1)
	case r.Int2 < 1:
		err = fmt.Errorf("Int2 must be >= 1, value %d was given", r.Int2)
	}
	return err
}

// Response represents a response that will be returned when a request is rendered
type Response struct {
	Lines chan string
}

// NewResponse is the Response factory
func NewResponse() *Response {
	return &Response{
		Lines: make(chan string),
	}
}
