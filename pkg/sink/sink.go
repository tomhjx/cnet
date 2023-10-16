package sink

import "github.com/tomhjx/cnet/pkg/core"

type Sink interface {
	Run(*string) error
}

type SinkInitialize interface {
	New(Option) (Sink, error)
}

type Option struct {
	Network core.Network `json:"network,omitempty"`
	Addr    string       `json:"addr,omitempty"`
}
