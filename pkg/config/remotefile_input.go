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

	c.ModTime = i.ParseModTime(resp)

	if c.Body, err = io.ReadAll(resp.Body); err != nil {
		return c, err
	}
	return c, nil
}

func (i *RemoteFileInput) ParseModTime(resp *http.Response) time.Time {
	modTime := time.Now()

	keys := map[string]string{
		"date":          http.TimeFormat,
		"last-modified": http.TimeFormat,
	}

	for k, format := range keys {
		vv := resp.Header.Get(k)
		if vv == "" {
			continue
		}
		mtime, err := time.Parse(format, vv)
		if err != nil {
			continue
		}
		return mtime
	}

	return modTime

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
