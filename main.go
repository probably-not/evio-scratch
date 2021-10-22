// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/probably-not/evio-scratch/internal"
	"github.com/tidwall/evio"
)

func main() {
	var port int
	var loops int

	flag.IntVar(&port, "port", 8080, "server port")
	flag.IntVar(&loops, "loops", 1, "num loops")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	handler := internal.NewHandler(ctx, loops, port)

	go func() {
		err := evio.Serve(handler, "tcp://127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	time.Sleep(time.Second * 10)
}
