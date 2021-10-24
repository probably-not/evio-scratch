//go:build !go1.16
// +build !go1.16

package evio

import (
	"io"
	"io/ioutil"
)

func closer(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func readall(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
