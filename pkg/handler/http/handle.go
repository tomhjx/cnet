package http

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/handler"
)

type Handle struct {
	Option *handler.Option
}

func (h Handle) Initialize(o *handler.Option) (handler.Handler, error) {
	return Handle{Option: o}, nil
}

func (h Handle) Do(hreq *core.CompletedRequest) (*core.Result, error) {
	if hreq.Method == "" {
		hreq.Method = http.MethodGet
	}
	req, err := http.NewRequest(hreq.Method, hreq.ADDR, nil)
	if err != nil {
		return nil, err
	}
	res := core.NewResult()
	var t0, t1, t2, t3, t4, t5, t6 time.Time
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t1 = time.Now()
		},
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				// connecting to IP
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Fatalf("unable to connect to host %v: %v", addr, err)
			}
			t2 = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
		GotFirstResponseByte: func() { t4 = time.Now() },
		TLSHandshakeStart:    func() { t5 = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t6 = time.Now() },
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// always refuse to follow redirects, visit does that
			// manually if required.
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Print SSL/TLS version which is used for connection
	res.RunTime.ConnectedVia = "plaintext"
	if resp.TLS != nil {
		switch resp.TLS.Version {
		case tls.VersionTLS12:
			res.RunTime.ConnectedVia = "TLSv1.2"
		case tls.VersionTLS13:
			res.RunTime.ConnectedVia = "TLSv1.3"
		}
	}
	if _, ok := h.Option.IncludeFields[field.Body]; ok {
		body, _ := h.readResponseBody(req, resp)
		res.Response.Body = &body
	}
	resp.Body.Close()
	res.Response.Status = resp.Status
	res.Response.StatusCode = resp.StatusCode

	if _, ok := h.Option.IncludeFields[field.Headers]; ok {
		for k, v := range resp.Header {
			res.Response.Headers[k] = strings.Join(v, ",")
		}
	}

	t7 := time.Now() // after read body
	if t0.IsZero() {
		// we skipped DNS
		t0 = t1
	}

	switch hreq.URL.Scheme {
	case "https":
		res.RunTime.AppConnectTime = t6.Sub(t0)
		res.RunTime.SSLTime = t6.Sub(t5)
	case "http":
		res.RunTime.AppConnectTime = t2.Sub(t0)
	}
	res.RunTime.NameLookUpTime = t1.Sub(t0)
	res.RunTime.ConnectTime = t2.Sub(t0)
	res.RunTime.TCPTime = t2.Sub(t1)
	res.RunTime.PreTransferTime = t3.Sub(t0)
	res.RunTime.StartTransferTime = t4.Sub(t0)
	res.RunTime.TTFB = res.RunTime.StartTransferTime - res.RunTime.AppConnectTime
	res.RunTime.ContentTransferTime = t7.Sub(t4)
	res.RunTime.ServerProcessTime = t4.Sub(t3)
	res.RunTime.TotalTime = t7.Sub(t0)
	return res, nil
}

func (h Handle) readResponseBody(req *http.Request, resp *http.Response) (string, error) {
	if h.isRedirect(resp) || req.Method == http.MethodHead {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func (h Handle) isRedirect(resp *http.Response) bool {
	return resp.StatusCode > 299 && resp.StatusCode < 400
}
