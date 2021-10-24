//go:build go1.16
// +build go1.16

package evio

import "io"

func closer(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

func readall(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
