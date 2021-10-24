// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"strconv"
	"time"

	cancellation "github.com/probably-not/evio-scratch/internal/cancellation"
	internalEvio "github.com/probably-not/evio-scratch/internal/evio"
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
	handler := internalEvio.NewHandler(ctx, loops, port)

	go func() {
		err := evio.Serve(handler, "tcp://127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	time.Sleep(time.Second * 10)
}
