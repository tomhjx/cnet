package flow

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/handler"
	"github.com/tomhjx/cnet/pkg/handler/host"
	"github.com/tomhjx/cnet/pkg/handler/http"
)

const (
	HTTPProtocol  Protocol = "http"
	HTTPsProtocol Protocol = "https"
	HOSTProtocol  Protocol = "host"
)

type ProtocolSet struct {
	handler     handler.Handler
	contentMeta []field.Field
}

var protocols = map[string]ProtocolSet{}

type Protocol string

func (p Protocol) Handler(o *handler.Option) (handler.Handler, error) {
	s, err := p.Load()
	if err != nil {
		return nil, err
	}
	return s.handler.Initialize(o)
}

func (p Protocol) Handle(req *core.CompletedRequest, o *handler.Option) (*core.Result, error) {
	h, err := p.Handler(o)
	if err != nil {
		return nil, err
	}
	req.ID = uuid.New().String()
	return h.Do(req)
}

func (p Protocol) ContentMeta() ([]field.Field, error) {
	s, err := p.Load()
	if err != nil {
		return nil, err
	}
	return s.contentMeta, nil
}

func (p Protocol) Load() (ProtocolSet, error) {
	protocol := string(p)
	if h, ok := protocols[protocol]; ok {
		return h, nil
	}
	return ProtocolSet{}, fmt.Errorf("not support protocol [%s]", protocol)
}

func (p Protocol) Inject(h handler.Handler, cm []field.Field) {
	protocols[string(p)] = ProtocolSet{handler: h, contentMeta: cm}
}

func ProtocolWithCompletedRequest(c *core.CompletedRequest) Protocol {
	if c.URL != nil {
		return Protocol(c.URL.Scheme)
	}
	return HOSTProtocol
}

func RegisterProtocols() {
	CommonMeta := []field.Field{field.ClientID, field.ADDR}
	httpMeta := append(CommonMeta, []field.Field{field.Method}...)
	HTTPProtocol.Inject(http.Handle{}, httpMeta)
	HTTPsProtocol.Inject(http.Handle{}, httpMeta)
	HOSTProtocol.Inject(host.Handle{}, CommonMeta)
}
