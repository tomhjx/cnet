package stdout

import (
	"fmt"
	"log"

	"github.com/tomhjx/cnet/pkg/sink"
)

type Sink struct {
	Option sink.Option
}

func (s Sink) New(o sink.Option) (sink.Sink, error) {
	log.Println("sink.stdout")
	return &Sink{Option: o}, nil
}

func (s *Sink) Run(c *string) error {
	fmt.Println(*c)
	return nil
}
