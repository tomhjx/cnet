package core

type Result struct {
	Response Response
	RunTime  RunTime
}

func NewResult() *Result {
	return &Result{
		Response: Response{
			Headers: map[string]string{},
		},
	}
}
