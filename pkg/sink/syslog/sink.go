package syslog

import (
	"fmt"
	"log/syslog"

	"github.com/tomhjx/cnet/pkg/sink"
)

type Sink struct {
	Option sink.Option
	syslog *syslog.Writer
}

func (s Sink) New(o sink.Option) (sink.Sink, error) {
	sl, err := syslog.Dial(string(o.Network), o.Addr, syslog.LOG_LOCAL1, "cnet")
	if err != nil {
		return nil, fmt.Errorf("network:%s,addr:%s,err:%s", o.Network, o.Addr, err)
	}
	return &Sink{Option: o, syslog: sl}, nil
}

func (s *Sink) Run(c *string) error {
	s.syslog.Info(*c)
	return nil
}
