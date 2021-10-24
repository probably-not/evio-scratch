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
	"strconv"
	"time"

	"github.com/panjf2000/gnet"
	cancellation "github.com/probably-not/evio-scratch/internal/cancellation"
	internalHttp "github.com/probably-not/evio-scratch/internal/http"
	"github.com/probably-not/evio-scratch/internal/ioutil"
	internalEvio "github.com/probably-not/evio-scratch/internal/loop/evio"
	internalGnet "github.com/probably-not/evio-scratch/internal/loop/gnet"
	"github.com/tidwall/evio"
)

var (
	port, loops            int
	help, useEvio, useGnet bool
)

func init() {
	flag.IntVar(&port, "port", 8080, "server port")
	flag.IntVar(&loops, "loops", 1, "num loops")
	flag.BoolVar(&help, "help", false, "show help message")
	flag.BoolVar(&useEvio, "evio", true, "use the evio event loop")
	flag.BoolVar(&useGnet, "gnet", false, "use the gnet event loop")
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(1)
	}

	if useEvio && useGnet {
		fmt.Println("multiple event loops specified, please use only one of evio or gnet")
		flag.Usage()
		os.Exit(2)
	}

	if !useEvio && !useGnet {
		fmt.Println("no event loops specified, please use one of evio or gnet")
		flag.Usage()
		os.Exit(2)
	}

	ctx := cancellation.CreateCancelContext()

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", internalHttp.Echo)

	if useEvio {
		handler := internalEvio.NewEvioLoop(ctx, loops, port, mux)

		go func() {
			err := evio.Serve(handler, "tcp://127.0.0.1:"+strconv.Itoa(port))
			if err != nil {
				panic(err)
			}
		}()
	} else if useGnet {
		handler := internalGnet.NewGnetLoop(ctx, loops, port, mux)

		go func() {
			err := gnet.Serve(handler, "tcp://127.0.0.1:"+strconv.Itoa(port), gnet.WithNumEventLoop(loops), gnet.WithLoadBalancing(gnet.RoundRobin))
			if err != nil {
				panic(err)
			}
		}()
	}

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
