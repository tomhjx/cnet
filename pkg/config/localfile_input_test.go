package config

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLocalFileInput(t *testing.T) {

	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}
	defer os.Remove(file.Name())

	var writeContentBody = func() []byte {
		b := []byte(uuid.New().String())
		os.WriteFile(file.Name(), b, 0644)
		return b
	}

	io, err := NewLocalFileInput(InputOption{Path: file.Name()})
	assert.Empty(t, err)

	wantBodys := sync.Map{}

	preContent := Content{ModTime: time.Now()}
	watched := make(chan struct{}, 10)
	go func() {
		io.Watch(func(i Inputer) {

			defer func() {
				watched <- struct{}{}
			}()

			t.Log("read when watch.")
			c, err := i.Read()
			assert.Empty(t, err)
			assert.True(t, c.ModTime.After(preContent.ModTime), fmt.Sprintln(c.ModTime, preContent.ModTime))
			preContent = c

			wkey := string(c.Body)
			if _, ok := wantBodys.Load(wkey); !ok {
				return
			}
			wantBodys.Delete(wkey)
			t.Log(wkey, "matching.")
		})
	}()

	lastWantBody := []byte{}
	n := 3
	for i := 0; i < n; i++ {
		lastWantBody = writeContentBody()
		wantBodys.Store(string(lastWantBody), true)
		time.Sleep(time.Second)
	}
	lastwkey := string(lastWantBody)

	exist := true
	for i := 0; i < n; i++ {
		t.Log("wait watched.")
		<-watched
		if _, ok := wantBodys.Load(lastwkey); !ok {
			exist = false
			break
		}
	}
	assert.False(t, exist)
	io.Close()

}
