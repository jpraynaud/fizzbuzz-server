// Package render handles requests that will be rendered according to the FizzBuzz algorithm
package render

import (
	"reflect"
	"testing"
)

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
		{"Limit == 0", fields{0, 3, 5, "A", "B"}, true},
		{"Int1 < 1", fields{20, 0, 5, "A", "B"}, true},
		{"Int2 < 1", fields{20, 3, -5, "A", "B"}, true},
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

func TestRenderer_Render(t *testing.T) {
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
		want    []string
		wantErr bool
	}{
		{"Limit < 0", fields{-20, 3, 5, "A", "B"}, []string{}, true},
		{"Limit == 0", fields{0, 3, 5, "A", "B"}, []string{}, true},
		{"Int1 < 1", fields{20, 0, 5, "A", "B"}, []string{}, true},
		{"Int2 < 1", fields{20, 3, -5, "A", "B"}, []string{}, true},
		{"Str1 is empty", fields{10, 2, 3, "", "B"}, []string{"1", "", "B", "", "5", "B", "7", "", "B", ""}, false},
		{"Str1 and Str2 are empty", fields{10, 2, 3, "", ""}, []string{"1", "", "", "", "5", "", "7", "", "", ""}, false},
		{"Int1 == Int2", fields{10, 3, 3, "A", "B"}, []string{"1", "2", "AB", "4", "5", "AB", "7", "8", "AB", "10"}, false},
		{"Limit < Int1", fields{10, 30, 3, "A", "B"}, []string{"1", "2", "B", "4", "5", "B", "7", "8", "B", "10"}, false},
		{"Limit < Int2", fields{10, 3, 30, "A", "B"}, []string{"1", "2", "A", "4", "5", "A", "7", "8", "A", "10"}, false},
		{"Int2 > Int1", fields{20, 5, 3, "A", "B"}, []string{"1", "2", "B", "4", "A", "B", "7", "8", "B", "A", "11", "B", "13", "14", "AB", "16", "17", "B", "19", "A"}, false},
		{"Standard case", fields{20, 3, 5, "A", "B"}, []string{"1", "2", "A", "4", "B", "A", "7", "8", "A", "B", "11", "A", "13", "14", "AB", "16", "17", "A", "19", "B"}, false},
		{"Unicode case", fields{20, 3, 5, "喂", "世界"}, []string{"1", "2", "喂", "4", "世界", "喂", "7", "8", "喂", "世界", "11", "喂", "13", "14", "喂世界", "16", "17", "喂", "19", "世界"}, false},
	}
	renderer := NewRenderer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &Request{
				Limit: tt.fields.Limit,
				Int1:  tt.fields.Int1,
				Int2:  tt.fields.Int2,
				Str1:  tt.fields.Str1,
				Str2:  tt.fields.Str2,
			}
			got := make([]string, 0)
			response := renderer.Render(request)
			for line := range response.Lines {
				got = append(got, line)
			}
			if (response.Error != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", response.Error, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Renderer.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}