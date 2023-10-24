package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/handler"
)

func TestHandle_Do(t *testing.T) {
	type args struct {
		hreq   *core.Request
		option *handler.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			"http",
			args{
				hreq: &core.Request{ADDR: "http://www.baidu.com/", Method: http.MethodGet},
			},
			nil,
		},
		{
			"https",
			args{
				hreq: &core.Request{ADDR: "https://www.baidu.com/", Method: http.MethodGet},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := Handle{}.Initialize(tt.args.option)
			assert.Empty(t, err)
			tt.args.hreq.Complete()
			res, err := h.Do(tt.args.hreq.Complete())
			assert.Equal(t, tt.wantErr, err)
			assert.NotEmpty(t, res.Response.Body)
			assert.NotEmpty(t, res.Response.Headers)
			assert.True(t, res.RunTime.ConnectTime > 0)

		})
	}
}
