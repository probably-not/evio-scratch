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
type responseWriter struct {
	*http.Response
	buf []byte
}

func NewResponseWriter() *responseWriter {
	return &responseWriter{
		Response: &http.Response{
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		},
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.Response.Header
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.StatusCode == 0 {
		rw.WriteHeader(200)
	}

	rw.buf = append(rw.buf, data...)
	return len(data), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func (rw *responseWriter) WriteToBuf(w io.Writer) error {
	rw.Body = ioutil.NopCloser(bytes.NewReader(rw.buf))
	rw.ContentLength = int64(len(rw.buf))
	return rw.Response.Write(w)
}
