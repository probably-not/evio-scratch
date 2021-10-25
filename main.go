// Copyright 2021 Coby Benveniste. All rights reserved.
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
	"path"
	"sync"
	"time"

	cancellation "github.com/probably-not/server-scratch/internal/cancellation"
	internalHttp "github.com/probably-not/server-scratch/internal/http"
	"github.com/probably-not/server-scratch/internal/ioutil"
	"github.com/probably-not/server-scratch/internal/loop"
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

	if engineType < 1 || engineType > 8 || engineType == loop.UnknownEngineType {
		fmt.Println("unknown engine type specified")
		flag.Usage()
		os.Exit(2)
	}

	ctx := cancellation.CreateCancelContext()

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", internalHttp.Echo)
	mux.HandleFunc("/sleep", internalHttp.Sleep)

	server, err := loop.NewServer(ctx, engineType, port, loops, mux)
	if err != nil {
		panic(err)
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	// Sleep for 1 second to ensure the server has started up
	time.Sleep(time.Second)

	err = testServer(10, "/echo")
	if err != nil {
		fmt.Println("Error in testing echo endpoint")
		panic(err)
	}
	err = testServer(5, "/sleep")
	if err != nil {
		fmt.Println("Error in testing sleep endpoint")
		panic(err)
	}
	fmt.Println("Completed testing the server, waiting for signal")

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

func testServer(reqs int, endpoint string) error {
	fmt.Println("Starting server tests for", endpoint)
	url := "http://" + path.Join("127.0.0.1:8080/", endpoint)

	wg := sync.WaitGroup{}
	wg.Add(reqs)

	for i := 0; i < reqs; i++ {
		time.Sleep(time.Millisecond * 200)
		go func(j int) {
			defer wg.Done()

			body := fmt.Sprintf(`{"req": %d}`, j)
			resp, err := http.Post(url, "application/json", bytes.NewReader([]byte(body)))
			if err != nil {
				fmt.Printf("error POST-ing for request %d. Error: %v\n", j, err)
				return
			}

			r, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error reading response body for request %d. Error: %v\n", j, err)
				return
			}

			if !bytes.Equal([]byte(body), r) {
				fmt.Println("Received unequal bytes!!! On request", j)
			}
			fmt.Println("Sent:", body, "Received:", string(r), "Headers:", resp.Header)
		}(i)
	}

	wg.Wait()

	return nil
}
