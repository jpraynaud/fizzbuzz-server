// Package render handles requests that will be rendered according to the FizzBuzz algorithm
package render

import (
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	type args struct {
		limit int
		int1  int
		int2  int
		str1  string
		str2  string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"NewRequest 1", args{1, 2, 3, "A", "B"}, &Request{1, 2, 3, "A", "B"}},
		{"NewRequest 2", args{-1, -2, -3, "AA", "BB"}, &Request{-1, -2, -3, "AA", "BB"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(tt.args.limit, tt.args.int1, tt.args.int2, tt.args.str1, tt.args.str2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Validate(t *testing.T) {
	type fields struct {
		Limit int
		Int1  int
		Int2  int
		Str1  string
		Str2  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Limit < 0", fields{-20, 3, 5, "A", "B"}, true},
		{"Int1 < 1", fields{20, 0, 5, "A", "B"}, true},
		{"Int2 < 1", fields{20, 3, -5, "A", "B"}, true},
		{"Limit == 0", fields{0, 3, 5, "A", "B"}, false},
		{"Str1 is empty", fields{10, 2, 3, "", "B"}, false},
		{"Str1 and Str2 are empty", fields{10, 2, 3, "", ""}, false},
		{"Int1 == Int2", fields{10, 3, 3, "A", "B"}, false},
		{"Limit < Int1", fields{10, 30, 3, "A", "B"}, false},
		{"Limit < Int2", fields{10, 3, 30, "A", "B"}, false},
		{"Int2 > Int1", fields{20, 5, 3, "A", "B"}, false},
		{"Standard case", fields{20, 3, 5, "A", "B"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Limit: tt.fields.Limit,
				Int1:  tt.fields.Int1,
				Int2:  tt.fields.Int2,
				Str1:  tt.fields.Str1,
				Str2:  tt.fields.Str2,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
