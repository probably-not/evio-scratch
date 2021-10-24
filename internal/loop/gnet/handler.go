package gnet

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/panjf2000/gnet"
	internalHttp "github.com/probably-not/evio-scratch/internal/http"
	"github.com/probably-not/evio-scratch/internal/loop"
	"github.com/tidwall/evio"
)

type GnetHandler struct {
	ctx         context.Context
	loops, port int
	httpHandler http.Handler
	*gnet.EventServer
}

func NewGnetLoop(ctx context.Context, loops, port int, httpHandler http.Handler) *GnetHandler {
	handler := GnetHandler{
		ctx:         ctx,
		loops:       loops,
		port:        port,
		httpHandler: httpHandler,
		EventServer: &gnet.EventServer{},
	}

	return &handler
}

// OnInitComplete fires on server up (one time)
func (h *GnetHandler) OnInitComplete(server gnet.Server) gnet.Action {
	fmt.Println("gnet server started with", server.NumEventLoop, "event loops on address", h.port)

	select {
	case <-h.ctx.Done():
		return gnet.Shutdown
	default:
		return gnet.None
	}
}

// OnOpened fires on opening new connections (per connection)
func (h *GnetHandler) OnOpened(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(&evio.InputStream{})

	select {
	case <-h.ctx.Done():
		return nil, gnet.Close
	default:
		return nil, gnet.None
	}
}

// OnClosed fires on closing connections (per connection)
func (h *GnetHandler) OnClosed(c gnet.Conn, err error) gnet.Action {
	if err != nil {
		fmt.Println("connection between", c.LocalAddr(), "and", c.RemoteAddr(), "has been closed with error value", err)
	}

	select {
	case <-h.ctx.Done():
		return gnet.Shutdown
	default:
		return gnet.None
	}
}

// React fires on data being sent to a connection (per connection, per data frame read)
func (h *GnetHandler) React(in []byte, c gnet.Conn) ([]byte, gnet.Action) {
	if len(in) == 0 {
		return nil, gnet.None
	}

	stream := c.Context().(*evio.InputStream)
	data := stream.Begin(in)

	complete, err := loop.IsRequestComplete(data)
	if err != nil {
		fmt.Println("Uh oh, there was an error checking completeness?", err)
		return nil, gnet.Close
	}

	stream.End(data)
	if !complete {
		return nil, gnet.None
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))
	if err != nil {
		fmt.Println("Uh oh, there was an error creating the request?", err)
		return nil, gnet.Close
	}

	res := internalHttp.NewResponseWriter()
	h.httpHandler.ServeHTTP(res, req)

	buf := bytes.NewBuffer(nil)
	err = res.WriteToBuf(buf)
	if err != nil {
		fmt.Println("Uh oh, there was an error writing the response?", err)
		return nil, gnet.Close
	}

	select {
	case <-h.ctx.Done():
		return nil, gnet.Close
	default:
		// Reset the connection context to an empty input stream once we have completed a full request in order to
		// ensure that the next request starts empty.
		c.SetContext(&evio.InputStream{})
		return buf.Bytes(), gnet.None
	}
}

func (h *GnetHandler) Tick() (delay time.Duration, action gnet.Action) {
	select {
	case <-h.ctx.Done():
		return time.Second, gnet.Shutdown
	default:
		return time.Second, gnet.None
	}
}
