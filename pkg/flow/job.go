package flow

import (
	"context"

	"github.com/google/uuid"
)

type Job struct {
	id          string
	ctx         context.Context
	tasks       []*Task
	onCancelled func()
}

func (j *Job) Run() error {
	for _, t := range j.tasks {
		go func(t *Task) {
			t.RunLoop()
		}(t)
	}
	return nil
}

func (j *Job) Cancel() error {
	j.onCancelled()
	return nil
}

func (j *Job) ID() string {
	return j.id
}

func NewJob(ctx context.Context, taskCtxs []*TaskContext) (*Job, error) {
	jobCtx, cancel := context.WithCancel(ctx)
	j := &Job{
		id:          uuid.New().String(),
		ctx:         jobCtx,
		tasks:       []*Task{},
		onCancelled: func() { cancel() },
	}
	for _, tc := range taskCtxs {
		t := NewTask(jobCtx, tc.Request.Complete(), tc.Sinks, tc.Option)
		t.request.JobID = j.ID()
		j.tasks = append(j.tasks, t)
	}
	return j, nil
}
