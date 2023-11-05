package flow

import (
	"encoding/json"
	"fmt"

	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/xlog"
)

type RawContent struct {
	Request *core.CompletedRequest
	Result  *core.Result
}

func NewFormatContent(rc *RawContent, fields []field.Field) *string {
	c := map[string]interface{}{}
	for _, f := range fields {
		v, err := f.ValueOf(rc)
		if err != nil {
			xlog.ErrorS(err, fmt.Sprintf("field [%s] format fail", f))
			continue
		}
		c[string(f)] = v
	}
	b, _ := json.Marshal(c)
	s := string(b)
	return &s
}
