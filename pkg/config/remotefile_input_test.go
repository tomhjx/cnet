package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoteFileInput(t *testing.T) {

	c1 := Content{
		Body:    []byte("Hello, World!"),
		ModTime: time.Now().UTC().Truncate(time.Second),
	}

	// 创建一个模拟的 HTTP 服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头
		w.Header().Set("date", time.Now().UTC().Format(http.TimeFormat))
		// 写入响应内容
		w.Write(c1.Body)
	}))
	defer ts.Close()

	in, err := NewRemoteFileInput(InputOption{Path: ts.URL, PollInterval: time.Second})
	assert.Empty(t, err)

	watched := make(chan struct{})

	go func() {
		in.Watch(func(i Inputer) {
			c, err := i.Read()
			t.Log(c)
			assert.Empty(t, err)
			assert.Equal(t, c1.Body, c.Body)
			assert.True(t, c.ModTime.After(c1.ModTime))
			c1 = c
			watched <- struct{}{}
		})
	}()

	for i := 0; i < 10; i++ {
		<-watched
	}
	in.Close()

}
