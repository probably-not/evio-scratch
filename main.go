// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	cancellation "github.com/probably-not/evio-scratch/internal/cancellation"
	internalHttp "github.com/probably-not/evio-scratch/internal/http"
	"github.com/probably-not/evio-scratch/internal/ioutil"
	"github.com/probably-not/evio-scratch/internal/loop"
)

var (
	port, loops int
	help        bool
	engineType  loop.EngineType
)

func init() {
	flag.IntVar(&port, "port", 8080, "server port")
	flag.IntVar(&loops, "loops", 1, "num loops")
	flag.BoolVar(&help, "help", false, "show help message")
}

func main() {
	flag.Var(&engineType, "engine", "engine type to use; can be one of stdlib, evio, or gnet")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(1)
	}

	if engineType == loop.UnknownEngineType {
		fmt.Println("unknown engine type specified")
		flag.Usage()
		os.Exit(2)
	}

	ctx := cancellation.CreateCancelContext()

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", internalHttp.Echo)

	server := loop.NewServer(ctx, engineType, port, loops, mux)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	testServer(10)

	<-ctx.Done()
	fmt.Println("Received exit signal, waiting 5 seconds to close gracefully")

	i := 0
	for range time.Tick(time.Second) {
		fmt.Print(".")
		i++
		if i >= 5 {
			os.Exit(0)
		}
	}
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
