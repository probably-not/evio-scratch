package http

import (
	"strconv"
	"testing"
)

var (
	complete    bool
	completeErr error
)

// BenchmarkParser_IsRequestComplete-10    	15606697	        77.08 ns/op	       0 B/op	       0 allocs/op
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
// BenchmarkParser_Strconv-10    	74064654	        16.22 ns/op	       7 B/op	       0 allocs/op
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
// BenchmarkParser_ParseContentLength-10    	125843733	         9.409 ns/op	       0 B/op	       0 allocs/op (switch table)
// BenchmarkParser_ParseContentLength-10    	158928097	         7.401 ns/op	       0 B/op	       0 allocs/op (slice table)
// --------------------------------------------------------------------------------------------------------------------------
// Benchmarks of math.Pow vs Pow lookup table:
// BenchmarkParser_ParseContentLength-10    	155882090	         7.545 ns/op	       0 B/op	       0 allocs/op (math.Pow)
// BenchmarkParser_ParseContentLength-10    	157663156	         7.404 ns/op	       0 B/op	       0 allocs/op (table)
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
