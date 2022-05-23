package http

/*
----------------------------------------------------------------------------------------------------
Testing and Benchmarking Cases for `IsRequestComplete(data []byte) (bool, error)`

The example request to be parsed is very simple and looks like this:
```
POST /echo HTTP/1.1
Host: 127.0.0.1:8080
User-Agent: Go-http-client/1.1
Content-Length: 10
Content-Type: application/json
Accept-Encoding: gzip

{"req": 0}
```
----------------------------------------------------------------------------------------------------
*/
var isRequestCompleteTestCases = []struct {
	expectedErr error
	desc        string
	input       []byte
	expected    bool
	wantErr     bool
}{
	{
		desc:        "incomplete headers",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\n"),
		expected:    false,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "complete headers empty body",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n"),
		expected:    true,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "complete headers non-empty body no content length",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n{\"req\": 0}"),
		expected:    false,
		wantErr:     true,
		expectedErr: errBadRequest,
	},
	{
		desc:        "complete headers with content length no body yet",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 10\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n"),
		expected:    false,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "complete headers with content length and incomplete body",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 10\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"req\": "),
		expected:    false,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "complete headers with content length and complete body",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 10\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"req\": 0}"),
		expected:    true,
		wantErr:     false,
		expectedErr: nil,
	},
	{
		desc:        "complete headers with bad content length",
		input:       []byte("POST /echo HTTP/1.1\r\nHost: 127.0.0.1:8080\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 123abc\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"req\": 0}"),
		expected:    false,
		wantErr:     true,
		expectedErr: errBadRequest,
	},
}

/*
----------------------------------------------------------------------------------------------------
Testing and Benchmarking Cases for `parseContentLength(clen []byte) (int64, error)`
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
