//go:build go1.16
// +build go1.16

package ioutil

import "io"

func NopCloser(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
