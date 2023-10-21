package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ConfigDemoOfTest struct {
	Name     string
	Desc     string
	Interval time.Duration
}

func TestWatch(t *testing.T) {

	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}
	var writeTestDemo = func() ConfigDemoOfTest {
		dd := ConfigDemoOfTest{}
		dd.Name = uuid.New().String()
		dd.Desc = fmt.Sprintf("desc %d", time.Now().UnixMicro())
		dd.Interval = time.Minute
		s, _ := json.Marshal(dd)
		os.WriteFile(file.Name(), s, 0644)
		return dd
	}

	// change
	waitings := sync.Map{}
	dd := writeTestDemo()
	waitings.Store(dd.Name, dd)
	defer os.Remove(file.Name())

	d := &ConfigDemoOfTest{}
	b, err := NewBuilder(InputOption{Path: file.Name()}, d)
	if err != nil {
		t.Fatal(err)
	}

	loaded := make(chan bool, 100)

	b.OnLoad(func(data any) {
		defer func() {
			loaded <- true
		}()
		d1 := data.(*ConfigDemoOfTest)
		d2any, ok := waitings.Load(d1.Name)
		if !ok {
			return
		}
		d2 := d2any.(ConfigDemoOfTest)
		assert.Equal(t, d2.Name, d1.Name)
		assert.Equal(t, d2.Desc, d1.Desc)
		assert.Equal(t, d2.Interval, d1.Interval)
		waitings.Delete(d2.Name)
		t.Log(d1.Name, "matching.")
	})

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	b.Watch()

	finalConfig := ConfigDemoOfTest{}
	n := 10
	for i := 0; i < n; i++ {
		finalConfig = writeTestDemo()
	}
	waitings.Store(finalConfig.Name, finalConfig)

	t.Log("waitings completed.")

	exist := false
	for i := 0; i < n; i++ {
		<-loaded
		exist = false
		waitings.Range(func(key, value any) bool {
			exist = true
			return false
		})
		if !exist {
			break
		}
	}
	assert.Equal(t, false, exist)
}
