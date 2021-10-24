package http

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

/*
----------------------------------------------------------------------------------------------------
Testing and Benchmarking IsRequestComplete(data []byte) (bool, error)
----------------------------------------------------------------------------------------------------
*/

// TODO: Tests!
func TestParser_IsRequestComplete(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

/*
----------------------------------------------------------------------------------------------------
Testing and Benchmarking parseContentLength(clen []byte) (int64, error)
----------------------------------------------------------------------------------------------------
*/

var parseContentLengthTestCases = []struct {
	expectedErr error
	desc        string
	input       []byte
	expected    int64
	wantErr     bool
}{
	{
		desc:        "first byte error",
		input:       []byte("a"),
		expected:    -1,
		wantErr:     true,
		expectedErr: errBadRequest,
	},
	{
		desc:        "middle byte error",
		input:       []byte("12a"),
		expected:    -1,
		wantErr:     true,
		expectedErr: errBadRequest,
	},
	{
		desc:        "0",
		input:       []byte("0"),
		expected:    0,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "1",
		input:       []byte("1"),
		expected:    1,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "2",
		input:       []byte("2"),
		expected:    2,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "3",
		input:       []byte("3"),
		expected:    3,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "4",
		input:       []byte("4"),
		expected:    4,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "5",
		input:       []byte("5"),
		expected:    5,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "6",
		input:       []byte("6"),
		expected:    6,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "7",
		input:       []byte("7"),
		expected:    7,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "8",
		input:       []byte("8"),
		expected:    8,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "9",
		input:       []byte("9"),
		expected:    9,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "10",
		input:       []byte("10"),
		expected:    10,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "123",
		input:       []byte("123"),
		expected:    123,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "1234",
		input:       []byte("1234"),
		expected:    1234,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "12345",
		input:       []byte("12345"),
		expected:    12345,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "123456",
		input:       []byte("123456"),
		expected:    123456,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "023456",
		input:       []byte("023456"),
		expected:    -1,
		wantErr:     true,
		expectedErr: errBadRequest,
	},
}

func TestParser_ParseContentLength(t *testing.T) {
	for _, tC := range parseContentLengthTestCases {
		t.Run(tC.desc, func(subT *testing.T) {
			got, err := parseContentLength(tC.input)
			if (err != nil) != tC.wantErr {
				subT.Errorf("parseContentLength() error = %v, wantErr %v", err, tC.wantErr)
				return
			}

			if err != nil && err != tC.expectedErr {
				subT.Errorf("parseContentLength() error type mismatch expecting %s and got %s", tC.expectedErr.Error(), err.Error())
				return
			}

			if !reflect.DeepEqual(got, tC.expected) {
				subT.Errorf("parseContentLength() got = %v, want %v", got, tC.expected)
			}
		})
	}

	// After the basic test cases let's run a bunch of tests on random numbers
	for i := 0; i < 100; i++ {
		r := rand.Int63n(1000000)
		b := []byte(fmt.Sprintf("%d", r))
		got, err := parseContentLength(b)
		if err != nil {
			t.Errorf("parseContentLength() error type mismatch expecting nil and got %s", err.Error())
			return
		}

		if !reflect.DeepEqual(got, r) {
			t.Errorf("parseContentLength() got = %v, want %v", got, r)
		}
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
