package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/probably-not/server-scratch/internal/ioutil"
)

// A very basic naive http.ResponseWriter implementation that attempts to write to an underlying http.Response.
// This should be further extended in the future to ensure we are writing the correct Headers, protocols, and flags
// to the http.Response.
type ResponseWriter struct {
	*http.Response
	buf []byte
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		Response: &http.Response{
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		},
	}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.Response.Header
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	if rw.StatusCode == 0 {
		rw.WriteHeader(200)
	}

	rw.buf = append(rw.buf, data...)
	return len(data), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func (rw *ResponseWriter) WriteToBuf(w io.Writer) error {
	rw.Body = ioutil.NopCloser(bytes.NewReader(rw.buf))
	rw.ContentLength = int64(len(rw.buf))
	return rw.Response.Write(w)
}
