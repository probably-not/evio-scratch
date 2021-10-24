package stdlib

import (
	"fmt"
	"net/http"
)

type Stdlib struct {
	*http.Server
}

func NewStdlib(port int, handler http.Handler) *Stdlib {
	fmt.Println("stdlib server started on address", port)

	return &Stdlib{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}
