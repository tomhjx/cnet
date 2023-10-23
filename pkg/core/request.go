package core

import (
	"net/url"
)

type Request struct {
	ID       string
	JobID    string
	TaskID   string
	ClientID string
	Tags     map[string]string
	RawURL   string
	Method   string
	Host     string
}

func (r *Request) Complete() *CompletedRequest {
	cr := &CompletedRequest{Request: r}
	u, err := r.CompleteURL()
	if err == nil {
		cr.URL = u
	}
	return cr
}

func (r *Request) CompleteURL() (*url.URL, error) {
	if r.RawURL == "" {
		return nil, nil
	}
	return url.Parse(r.RawURL)
}

type CompletedRequest struct {
	*Request
	URL *url.URL `json:"-"`
}
