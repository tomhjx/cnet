package flow

import (
	"encoding/json"

	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
)

type RawContent struct {
	Request *core.CompletedRequest
	Result  *core.Result
}

func NewFormatContent(rc *RawContent, fields []field.Field) *string {
	c := map[string]interface{}{}
	for _, f := range fields {
		c[string(f)] = f.ValueOf(rc)
	}
	b, _ := json.Marshal(c)
	s := string(b)
	return &s
}
