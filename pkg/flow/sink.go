package flow

import (
	"fmt"

	"github.com/tomhjx/cnet/pkg/sink"
	"github.com/tomhjx/cnet/pkg/sink/stdout"
	"github.com/tomhjx/cnet/pkg/sink/syslog"
)

const (
	StdOutSinkName = "stdout"
	SysLogSinkName = "syslog"
)

var sinkInitializes = map[string]sink.SinkInitialize{}
var sinks = map[SinkIndex]sink.Sink{}

type SinkIndex struct {
	Name   string
	Option sink.Option
}

func RegisterSinks() {
	sinkInitializes[StdOutSinkName] = stdout.Sink{}
	sinkInitializes[SysLogSinkName] = syslog.Sink{}
}

func NewSink(name string, o sink.Option) (sink.Sink, error) {
	if s, ok := sinkInitializes[name]; ok {
		return s.New(o)
	}
	return nil, fmt.Errorf("not support sink [%s]", name)
}
func SinkOf(name string, o sink.Option) (sink.Sink, error) {
	si := SinkIndex{
		Name:   name,
		Option: o,
	}
	if s, ok := sinks[si]; ok {
		return s, nil
	}
	s, err := NewSink(name, o)
	if err != nil {
		return s, err
	}
	sinks[si] = s
	return s, nil
}

type SinkConfig struct {
	Name   string      `json:"name"`
	Option sink.Option `json:"option,omitempty"`
}
