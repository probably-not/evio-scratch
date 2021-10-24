package loop

type Engine interface {
	ListenAndServe() error
}

type EngineType uint32

const (
	Stdlib EngineType = 1 << iota
	Evio
	Gnet
	ErrUnknownEngineType
)

func (et EngineType) String() string {
	switch et {
	case Stdlib:
		return "Stdlib"
	case Evio:
		return "Evio"
	case Gnet:
		return "Gnet"
	case ErrUnknownEngineType:
		return "ErrUnknownEngineType"
	default:
		return ""
	}
}
