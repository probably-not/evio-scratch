//go:build !go1.16
// +build !go1.16

package ioutil

import (
	"io"
	"io/ioutil"
)

func NopCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
