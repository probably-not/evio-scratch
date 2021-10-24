// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cancellation "github.com/probably-not/evio-scratch/internal/cancellation"
	internalEvio "github.com/probably-not/evio-scratch/internal/evio"
	internalHttp "github.com/probably-not/evio-scratch/internal/http"
	"github.com/probably-not/evio-scratch/internal/ioutil"
	"github.com/tidwall/evio"
)

var port, loops int

func init() {
	flag.IntVar(&port, "port", 8080, "server port")
	flag.IntVar(&loops, "loops", 1, "num loops")
}

func main() {
	flag.Parse()

	ctx := cancellation.CreateCancelContext()

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", internalHttp.Echo)

	handler := internalEvio.NewHandler(ctx, loops, port, mux)

	go func() {
		err := evio.Serve(handler, "tcp://127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			panic(err)
		}
	}()

	testServer(10)

	<-ctx.Done()
	time.Sleep(time.Second * 5)
}

func testServer(reqs int) error {
	for i := 0; i < reqs; i++ {
		j := i
		body := fmt.Sprintf(`{"req": %d}`, j)
		resp, err := http.Post("http://127.0.0.1:8080/echo", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			return err
		}

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if !bytes.Equal([]byte(body), r) {
			fmt.Println("Received unequal bytes!!!")
		}
		fmt.Println("Sent:", body, "Received:", string(r))
	}

	return nil
}
