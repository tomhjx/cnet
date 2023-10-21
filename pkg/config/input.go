package config

import (
	"net/url"
	"path/filepath"
	"time"
)

type Inputer interface {
	Read() (Content, error)
	Watch(func(Inputer)) error
	Close()
}

type Content struct {
	Body    []byte
	ModTime time.Time
}

type InputOption struct {
	Path         string
	PollInterval time.Duration
}

func (o InputOption) NewInputer() (Inputer, error) {
	if filepath.IsAbs(o.Path) {
		return NewLocalFileInput(o)
	}

	if _, err := url.Parse(o.Path); err == nil {
		return NewRemoteFileInput(o)
	}

	return nil, nil
}
