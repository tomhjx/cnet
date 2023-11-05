package core

import (
	"net/url"
)

type Request struct {
	ID              string
	JobID           string
	TaskID          string
	ClientID        string
	Tags            map[string]string
	Method          string
	ADDR            string
	Queue           string
	QueueExchange   string
	QueueRoutingKey string
}

func (r *Request) Complete() *CompletedRequest {
	cr := &CompletedRequest{Request: r}
	u, _ := r.CompleteURL()
	if u.Scheme != "" {
		cr.URL = u
	}
	return cr
}

func (r *Request) CompleteURL() (*url.URL, error) {
	return url.Parse(r.ADDR)
}

type CompletedRequest struct {
	*Request
	URL *url.URL `json:"-"`
}
