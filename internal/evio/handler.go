package evio

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	internalHttp "github.com/probably-not/evio-scratch/internal/http"
	"github.com/tidwall/evio"
)

func NewHandler(ctx context.Context, loops, port int) evio.Events {
	var handler evio.Events
	handler.NumLoops = loops
	handler.LoadBalance = evio.RoundRobin

	// Serving fires on server up (one time)
	handler.Serving = func(server evio.Server) evio.Action {
		fmt.Println("evio server started with", server.NumLoops, "event loops on port", port)

		select {
		case <-ctx.Done():
			return evio.Shutdown
		default:
			return evio.None
		}
	}

	// Opened fires on opening new connections (per connection)
	handler.Opened = func(c evio.Conn) ([]byte, evio.Options, evio.Action) {
		c.SetContext(&evio.InputStream{})

		select {
		case <-ctx.Done():
			return nil, evio.Options{}, evio.Close
		default:
			return nil, evio.Options{}, evio.None
		}
	}

	// Closed fires on closing connections (per connection)
	handler.Closed = func(c evio.Conn, err error) evio.Action {
		if err != nil {
			fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "has been closed with error value", err)
		}

		select {
		case <-ctx.Done():
			return evio.Shutdown
		default:
			return evio.None
		}
	}

	// Data fires on data being sent to a connection (per connection, per data frame read)
	handler.Data = func(c evio.Conn, in []byte) ([]byte, evio.Action) {
		if len(in) == 0 {
			return nil, evio.None
		}

		stream := c.Context().(*evio.InputStream)
		data := stream.Begin(in)

		complete, err := isRequestComplete(data)
		if err != nil {
			fmt.Println("Uh oh, there was an error checking completeness?", err)
			return nil, evio.Close
		}

		stream.End(data)
		if !complete {
			return nil, evio.None
		}

		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))
		if err != nil {
			fmt.Println("Uh oh, there was an error creating the request?", err)
			return nil, evio.Close
		}

		res := NewResponseWriter()
		internalHttp.Echo(res, req)

		buf := bytes.NewBuffer(nil)
		err = res.WriteToBuf(buf)
		if err != nil {
			fmt.Println("Uh oh, there was an error writing the response?", err)
			return nil, evio.Close
		}

		select {
		case <-ctx.Done():
			return nil, evio.Close
		default:
			// Reset the connection context to an empty input stream once we have completed a full request in order to
			// ensure that the next request starts empty.
			c.SetContext(&evio.InputStream{})
			return buf.Bytes(), evio.None
		}
	}

	handler.Tick = func() (delay time.Duration, action evio.Action) {
		select {
		case <-ctx.Done():
			return time.Second, evio.Shutdown
		default:
			return time.Second, evio.None
		}
	}

	return handler
}
