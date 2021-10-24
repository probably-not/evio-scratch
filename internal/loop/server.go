package loop

import (
	"context"
	"net/http"

	"github.com/probably-not/evio-scratch/internal/loop/evio"
	"github.com/probably-not/evio-scratch/internal/loop/gnet"
	"github.com/probably-not/evio-scratch/internal/loop/stdlib"
)

type Server struct {
	ctx    context.Context
	engine Engine
}

func NewServer(ctx context.Context, engineType EngineType, port, loops int, handler http.Handler) *Server {
	var engine Engine
	switch engineType {
	case Evio:
		engine = evio.NewEngine(ctx, loops, port, handler)
	case Gnet:
		engine = gnet.NewEngine(ctx, loops, port, handler)
	default:
		engine = stdlib.NewStdlib(port, handler)
	}

	return &Server{
		ctx:    ctx,
		engine: engine,
	}
}
