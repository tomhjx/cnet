package main

import (
	"github.com/tomhjx/cnet/cmd/cnet/app"
	"github.com/tomhjx/xlog"
)

func main() {

	if err := app.NewCommand().Execute(); err != nil {
		xlog.Fatal(err)
	}
}
