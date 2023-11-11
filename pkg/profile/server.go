package profile

import (
	"context"
	"fmt"

	"net/http"
	_ "net/http/pprof"

	"github.com/tomhjx/cnet/pkg/xlogging"
	"github.com/tomhjx/xlog"
)

func StartServer(port int, ctx context.Context) {
	if port <= 0 {
		return
	}

	go func() {
		xlogging.Changed().Infof("profile server(:%d) started.", port)
		xlog.Error(http.ListenAndServe(fmt.Sprint(":", port), nil), "profile http server start fail.")
	}()
}
