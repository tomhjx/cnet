package app

import (
	"context"

	"github.com/tomhjx/cnet/cmd/cnet/app/options"
	"github.com/tomhjx/cnet/pkg/flow"
	"github.com/tomhjx/cnet/pkg/metric"
	"github.com/tomhjx/xlog"
)

var (
	job          *flow.Job
	jobPreloaded = make(chan *jobContext, 10)
)

type jobContext struct {
	ctx    context.Context
	config options.ConfigData
}

func NotifyLoadJob(ctx context.Context, config options.ConfigData) error {
	jobPreloaded <- &jobContext{ctx: ctx, config: config}
	return nil
}

func Run(ctx context.Context, config *options.Config, stopCh <-chan struct{}) error {
	go func() {
		for {
			jc := <-jobPreloaded
			if err := runJob(jc.ctx, jc.config.Complete()); err != nil {
				xlog.Fatal(err)
			}
		}
	}()
	config.Init(func() {
		metric.StartServer(config.MetricsServerPort, ctx)
	}, func(c options.ConfigData) {
		NotifyLoadJob(context.Background(), c)
	})

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
