package options

import (
	"io"
	"os"
)

type Options struct {
	Arguments []string
	Config    *Config
	IOStreams
}

type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

func NewOptions() (*Options, error) {
	s := Options{
		Config:    NewConfig(),
		Arguments: os.Args,
		IOStreams: IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	}
	return &s, nil
}
