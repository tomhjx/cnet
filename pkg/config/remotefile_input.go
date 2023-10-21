package config

import (
	"io"
	"log"
	"net/http"
	"time"
)

type RemoteFileInput struct {
	option InputOption
	done   chan struct{}
}

func NewRemoteFileInput(option InputOption) (Inputer, error) {
	return &RemoteFileInput{option: option, done: make(chan struct{})}, nil
}

func (i *RemoteFileInput) Read() (Content, error) {
	c := Content{}
	resp, err := http.Get(i.option.Path)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	lastModified := resp.Header.Get("Last-Modified")
	if c.ModTime, err = time.Parse(http.TimeFormat, lastModified); err != nil {
		return c, err
	}
	if c.Body, err = io.ReadAll(resp.Body); err != nil {
		return c, err
	}
	return c, nil
}

func (i *RemoteFileInput) Watch(cb func(Inputer)) error {
	ticker := time.NewTicker(i.option.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-i.done:
			log.Println("Watcher Done,", i.option.Path)
			return nil
		}
		cb(i)
	}
}

func (i *RemoteFileInput) Close() {
	close(i.done)
}
