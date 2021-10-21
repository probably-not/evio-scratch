package internal

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/tidwall/evio"
)

func NewHandler(loops, port int) evio.Events {
	var handler evio.Events
	handler.NumLoops = loops
	handler.LoadBalance = evio.RoundRobin

	handler.Serving = func(server evio.Server) evio.Action {
		fmt.Println("evio server started with", server.NumLoops, "event loops on port", port)
		return evio.None
	}

	handler.Opened = func(c evio.Conn) ([]byte, evio.Options, evio.Action) {
		fmt.Println("new connection opened between", c.LocalAddr(), "and", c.RemoteAddr())
		c.SetContext(&evio.InputStream{})
		return nil, evio.Options{}, evio.None
	}

	handler.Closed = func(c evio.Conn, err error) evio.Action {
		fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "has been closed with error value", err)
		return evio.None
	}

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

		fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "sending data", buf.String())
		return buf.Bytes(), evio.None
	}

	return handler
}

/*
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		if in == nil {
			return
		}
		is := c.Context().(*evio.InputStream)
		data := is.Begin(in)
		if noparse && bytes.Contains(data, []byte("\r\n\r\n")) {
			// for testing minimal single packet request -> response.
			out = appendresp(nil, "200 OK", "", res)
			return
		}
		// process the pipeline
		var req request
		for {
			leftover, err := parsereq(data, &req)
			if err != nil {
				// bad thing happened
				out = appendresp(out, "500 Error", "", err.Error()+"\n")
				action = evio.Close
				break
			} else if len(leftover) == len(data) {
				// request not ready, yet
				break
			}
			// handle the request
			req.remoteAddr = c.RemoteAddr().String()
			out = appendhandle(out, &req)
			data = leftover
		}
		is.End(data)
		return
	}
*/
