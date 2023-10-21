package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type LocalFileInput struct {
	option InputOption
	done   chan struct{}
}

func NewLocalFileInput(option InputOption) (Inputer, error) {
	return &LocalFileInput{option: option, done: make(chan struct{})}, nil
}

func (i *LocalFileInput) Read() (Content, error) {
	c := Content{}
	info, err := os.Stat(i.option.Path)
	if err != nil {
		return c, err
	}
	c.ModTime = info.ModTime()
	if c.Body, err = os.ReadFile(i.option.Path); err != nil {
		return c, err
	}
	return c, nil
}

func (i *LocalFileInput) Watch(cb func(Inputer)) error {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := w.Add(i.option.Path); err != nil {
		return err
	}

	defer w.Close()
	for {
		select {
		case <-i.done:
			log.Println("Watcher Done,", i.option.Path)
			return nil
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return err
			}
			log.Println("ERROR: ", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return nil
			}

			// We just want to watch for file creation, so ignore everything
			// outside of Create and Write.
			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
				continue
			}
			cb(i)
		}
	}
}

func (i *LocalFileInput) Close() {
	close(i.done)
}
