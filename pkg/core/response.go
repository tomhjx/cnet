package core

type Response struct {
	Headers    map[string]string
	Body       *string
	Status     string
	StatusCode int
}
