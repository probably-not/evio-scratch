package loop

import (
	"errors"
	"strings"
)

type Engine interface {
	ListenAndServe() error
}

type EngineType uint32

const (
	Stdlib EngineType = 1 << iota
	Evio
	Gnet
	UnknownEngineType
)

var ErrUnknownEngineType = errors.New("unknown engine type")

func (et EngineType) String() string {
	switch et {
	case Stdlib:
		return "Stdlib"
	case Evio:
		return "Evio"
	case Gnet:
		return "Gnet"
	case UnknownEngineType:
		return "UnknownEngineType"
	default:
		return ""
	}
}

func (et *EngineType) Set(value string) error {
	switch strings.ToLower(value) {
	case "evio":
		*et = Evio
	case "gnet":
		*et = Gnet
	case "stdlib":
		*et = Stdlib
	default:
		*et = UnknownEngineType
	}

	if *et == UnknownEngineType {
		return ErrUnknownEngineType
	}

	return nil
}
