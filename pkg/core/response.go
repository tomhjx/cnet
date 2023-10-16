package core

type Response struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    *string           `json:"body,omitempty"`
}
