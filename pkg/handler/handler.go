package handler

import (
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
)

type Handler interface {
	Initialize(o *Option) (Handler, error)
	Do(req *core.CompletedRequest) (*core.Result, error)
}

type Option struct {
	IncludeFields map[field.Field]bool
}
