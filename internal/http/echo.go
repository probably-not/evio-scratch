package http

import (
	"net/http"

	"github.com/probably-not/evio-scratch/internal/ioutil"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("unable to read request body"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(b)
	w.WriteHeader(http.StatusOK)
}
