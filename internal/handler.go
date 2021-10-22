package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

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
			fmt.Println("handler.Serving context is closed, we are shutting down")
			return evio.Shutdown
		default:
			return evio.None
		}
	}

	// Opened fires on opening new connections (per connection)
	handler.Opened = func(c evio.Conn) ([]byte, evio.Options, evio.Action) {
		fmt.Println("new connection opened between", c.LocalAddr(), "and", c.RemoteAddr())
		c.SetContext(&evio.InputStream{})

		select {
		case <-ctx.Done():
			fmt.Println("handler.Opened context is closed, we are no longer accepting connections")
			return nil, evio.Options{}, evio.Close
		default:
			return nil, evio.Options{}, evio.None
		}
	}

	// Closed fires on closing connections (per connection)
	handler.Closed = func(c evio.Conn, err error) evio.Action {
		fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "has been closed with error value", err)

		select {
		case <-ctx.Done():
			fmt.Println("handler.Closed context is closed, we are no longer accepting connections")
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

		fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "received data", string(in))

		body := []byte("Hello there")

		res := http.Response{
			StatusCode:    200,
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(body)),
			Close:         false,
			Body:          closer(bytes.NewReader(body)),
		}
		buf := bytes.NewBuffer(nil)
		err := res.Write(buf)
		if err != nil {
			fmt.Println("Uh oh, there was an error?", err)
			return nil, evio.Close
		}

		select {
		case <-ctx.Done():
			fmt.Println("handler.Data context is closed, we are no longer accepting connections")
			return nil, evio.Close
		default:
			fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "sending data", buf.String())
			return buf.Bytes(), evio.None
		}
	}

	handler.Tick = func() (delay time.Duration, action evio.Action) {
		select {
		case <-ctx.Done():
			fmt.Println("handler.Tick context is closed, we are no longer accepting connections")
			return time.Second, evio.Shutdown
		default:
			fmt.Println("handler.Tick")
			return time.Second, evio.None
		}
	}

	return handler
}
