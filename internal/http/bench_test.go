package http

import (
	"strconv"
	"testing"
)

var (
	complete    bool
	completeErr error
)

// Initial Benchmark Scores to beat:
// BenchmarkParser_IsRequestComplete-16    	13264154	        87.29 ns/op	       0 B/op	       0 allocs/op
func BenchmarkParser_IsRequestComplete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := IsRequestComplete(isRequestCompleteTestCases[i%len(isRequestCompleteTestCases)].input)
		if err != nil {
			completeErr = err
		}
		complete = c
	}
}

var (
	length   int64
	parseErr error
)

// Benchmark to test strconv vs a custom parsing function for parsing out content lengths.
// BenchmarkParser_Strconv-16    	43783546	        28.33 ns/op	       7 B/op	       0 allocs/op
func BenchmarkParser_Strconv(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, err := strconv.ParseInt(string(parseContentLengthTestCases[i%len(parseContentLengthTestCases)].input), 10, 64)
		if err != nil {
			parseErr = err
		}
		length = l
	}
}

// Benchmarks to find best Table structure:
// BenchmarkParser_ParseContentLength-16    	74605522	        14.63 ns/op	       0 B/op	       0 allocs/op (switch table)
// BenchmarkParser_ParseContentLength-16    	76251918	        13.59 ns/op	       0 B/op	       0 allocs/op (slice table)
// BenchmarkParser_ParseContentLength-16    	34651278	        32.38 ns/op	       0 B/op	       0 allocs/op (map table)
// --------------------------------------------------------------------------------------------------------------------------
// Benchmarks of math.Pow vs Pow lookup table:
// BenchmarkParser_ParseContentLength-16    	82135372	        13.56 ns/op	       0 B/op	       0 allocs/op (math.Pow)
// BenchmarkParser_ParseContentLength-16    	91160745	        12.85 ns/op	       0 B/op	       0 allocs/op (table)
func BenchmarkParser_ParseContentLength(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, err := parseContentLength(parseContentLengthTestCases[i%len(parseContentLengthTestCases)].input)
		if err != nil {
			parseErr = err
		}
		length = l
	}
}
