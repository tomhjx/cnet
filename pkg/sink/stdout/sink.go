package stdout

import (
	"fmt"

	"github.com/tomhjx/cnet/pkg/sink"
)

type Sink struct {
	Option sink.Option
}

func (s Sink) New(o sink.Option) (sink.Sink, error) {
	return &Sink{Option: o}, nil
}

func (s *Sink) Run(c *string) error {
	fmt.Println(*c)
	return nil
}
