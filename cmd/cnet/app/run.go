package app

import (
	"context"
	"log"

	"github.com/tomhjx/cnet/cmd/cnet/app/options"
	"github.com/tomhjx/cnet/pkg/flow"
)

var (
	job          *flow.Job
	jobPreloaded = make(chan *jobContext)
)

type jobContext struct {
	ctx    context.Context
	config *options.Config
}

func NotifyLoadJob(ctx context.Context, config *options.Config) error {
	jobPreloaded <- &jobContext{ctx: ctx, config: config}
	return nil
}

func Run(ctx context.Context, config *options.Config, stopCh <-chan struct{}) error {
	go func() {
		for {
			jc := <-jobPreloaded
			if err := runJob(jc.ctx, jc.config.Complete()); err != nil {
				log.Fatalln(err)
			}
		}
	}()
	NotifyLoadJob(ctx, config)
	<-stopCh
	return nil
}

func runJob(ctx context.Context, configs []*options.CompletedConfig) error {
	if job != nil {
		job.Cancel()
	}
	if err := loadJob(ctx, configs); err != nil {
		return err
	}
	return job.Run()
}

func loadJob(ctx context.Context, configs []*options.CompletedConfig) error {
	taskCtxs := []*flow.TaskContext{}
	for _, c := range configs {
		option := &flow.TaskOption{
			Interval:      c.Interval,
			IncludeFields: c.IncludeFields,
		}
		taskCtx, err := flow.NewTaskContext(option, c.CreateRequest(), c.Sinks)
		if err != nil {
			return err
		}
		taskCtxs = append(taskCtxs, taskCtx)
	}
	j, err := flow.NewJob(ctx, taskCtxs)
	if err != nil {
		return err
	}
	job = j
	return nil
}
