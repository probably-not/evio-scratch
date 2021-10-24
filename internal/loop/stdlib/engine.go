package stdlib

import (
	"fmt"
	"net/http"
)

type Stdlib struct {
	*http.Server
}

func NewStdlib(port int, handler http.Handler) *Stdlib {
	return &Stdlib{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}
