package main

import (
	"log"

	"github.com/tomhjx/cnet/cmd/cnet/app"
)

func main() {

	if err := app.NewCommand().Execute(); err != nil {
		log.Fatalf("err: %v", err)
	}
}
