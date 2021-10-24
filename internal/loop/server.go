package loop

import (
	"context"
	"net/http"

	"github.com/probably-not/server-scratch/internal/loop/evio"
	"github.com/probably-not/server-scratch/internal/loop/gnet"
	"github.com/probably-not/server-scratch/internal/loop/stdlib"
)

type Server struct {
	ctx    context.Context
	engine Engine
}

func NewServer(ctx context.Context, engineType EngineType, port, loops int, handler http.Handler) (*Server, error) {
	var engine Engine
	switch engineType {
	case Evio:
		engine = evio.NewEngine(ctx, loops, port, handler)
	case Gnet:
		engine = gnet.NewEngine(ctx, loops, port, handler)
	case Stdlib:
		engine = stdlib.NewStdlib(port, handler)
	case UnknownEngineType:
		return nil, ErrUnknownEngineType
	default:
		return nil, ErrUnknownEngineType
	}

	return &Server{
		ctx:    ctx,
		engine: engine,
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.engine.ListenAndServe()
}
