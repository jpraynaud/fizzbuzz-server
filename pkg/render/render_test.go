// Package render handles requests that will be rendered according to the FizzBuzz algorithm
package render

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRequest_Validate(t *testing.T) {
	// Prepare tests data
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
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request := NewRequest(tt.fields.Limit, tt.fields.Int1, tt.fields.Int2, tt.fields.Str1, tt.fields.Str2)
			// Check request validation matches the one wanted
			if err := request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderer_Render(t *testing.T) {
	// Prepare tests data
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
	// Create renderer
	renderer := NewRenderer()
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request := NewRequest(tt.fields.Limit, tt.fields.Int1, tt.fields.Int2, tt.fields.Str1, tt.fields.Str2)
			// Render request and convert to slice
			got := make([]string, 0)
			response := renderer.Render(context.TODO(), request)
			for item := range response.Items {
				got = append(got, item)
			}
			// Check that slice matches the one wanted
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

func BenchmarkRenderer_Render(b *testing.B) {
	// Create renderer
	renderer := NewRenderer()
	// Create request
	request := NewRequest(100, 3, 5, "fizz", "buzz")
	// Reset timer
	b.ResetTimer()
	// Run benchmark
	for i := 0; i < b.N; i++ {
		// Render request
		renderer.Render(context.TODO(), request)
	}
}

func ExampleRenderer_Render() {
	// Create renderer
	renderer := NewRenderer()
	// Create request
	request := NewRequest(20, 3, 5, "fizz", "buzz")
	// Render request
	response := renderer.Render(context.TODO(), request)
	// Generate output
	items := make([]string, 0)
	for item := range response.Items {
		items = append(items, item)
	}
	fmt.Print(strings.Join(items, ","))
	// Output: 1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz
}

func TestStatistics(t *testing.T) {
	// Prepare tests data
	type fields struct {
		Limit int
		Int1  int
		Int2  int
		Str1  string
		Str2  string
	}
	tests := []struct {
		name               string
		fieldsTodo         fields
		iterationsTodo     int
		iterationsWant     int
		fieldsMostUsedWant fields
	}{
		{"Step 1", fields{10, 3, 5, "A", "B"}, 1, 1, fields{10, 3, 5, "A", "B"}},
		{"Step 2", fields{20, 3, 5, "A", "B"}, 2, 2, fields{20, 3, 5, "A", "B"}},
		{"Step 3", fields{20, 3, 5, "A", "B"}, 10, 12, fields{20, 3, 5, "A", "B"}},
		{"Step 4", fields{20, 3, 5, "A", "B"}, 30, 42, fields{20, 3, 5, "A", "B"}},
		{"Step 5", fields{30, 3, 5, "A", "B"}, 100, 100, fields{30, 3, 5, "A", "B"}},
	}
	// Get & reset statistics
	statistics := NewStatistics()
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			request := NewRequest(tt.fieldsTodo.Limit, tt.fieldsTodo.Int1, tt.fieldsTodo.Int2, tt.fieldsTodo.Str1, tt.fieldsTodo.Str2)
			// Records request multiple times
			for i := 0; i < tt.iterationsTodo; i++ {
				statistics.RecordStatistic(request)
			}
			// Check that total renderings recorded matches the total wanted
			requestStatistic := statistics.GetStatistic(request)
			if got := requestStatistic.Total; got != tt.iterationsWant {
				t.Errorf("StatisticRecorder.RecordStatistic() = %v, want %v", got, tt.iterationsWant)
			}
			// Check that the most used request matches the most used request wanted
			topRequestWant := NewRequest(tt.fieldsMostUsedWant.Limit, tt.fieldsMostUsedWant.Int1, tt.fieldsMostUsedWant.Int2, tt.fieldsMostUsedWant.Str1, tt.fieldsMostUsedWant.Str2)
			topRequestStatistic := statistics.GetTopStatistic()
			if got := topRequestStatistic.Request; !reflect.DeepEqual(got, *topRequestWant) {
				t.Errorf("Renderer.Render() = %v, want %v", got, topRequestWant)
			}
		})
	}
}

func BenchmarkRenderer_GetTopStatistic(b *testing.B) {
	// Create renderer
	renderer := NewRenderer()
	// Create request
	request := NewRequest(100, 3, 5, "fizz", "buzz")
	// record statistics
	for i := 0; i < 1000; i++ {
		// Create request
		renderer.Render(context.TODO(), request)
	}
	// Reset timer
	b.ResetTimer()
	// Run benchmark
	for i := 0; i < b.N; i++ {
		// Create request
		renderer.GetTopStatistic()
	}
}
