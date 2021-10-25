package http

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/probably-not/server-scratch/internal/ioutil"
)

func Sleep(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("unable to read request body"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sleepTime := time.Second * time.Duration(rand.Intn(30))
	<-time.After(sleepTime)

	w.Header().Set("x-sleep-time", sleepTime.String())
	w.Write(b)
}
