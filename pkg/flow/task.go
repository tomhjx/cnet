package flow

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/handler"
	"github.com/tomhjx/cnet/pkg/sink"
)

type Task struct {
	id             string
	ctx            context.Context
	protocol       Protocol
	request        *core.CompletedRequest
	sinks          []sink.Sink
	option         *TaskOption
	handlerOption  *handler.Option
	interval       time.Duration
	remainingCount int64
	isAlways       bool
}

func NewTask(ctx context.Context, request *core.CompletedRequest, sinks []sink.Sink, option *TaskOption) *Task {
	ho := &handler.Option{
		IncludeFields: map[field.Field]bool{},
	}
	for _, v := range option.IncludeFields {
		ho.IncludeFields[v] = true
	}
	ho.TimeOut = option.Interval
	t := &Task{
		id:             uuid.New().String(),
		ctx:            ctx,
		protocol:       ProtocolWithCompletedRequest(request),
		request:        request,
		sinks:          []sink.Sink{},
		option:         option,
		handlerOption:  ho,
		interval:       option.Interval,
		remainingCount: option.Count,
	}
	t.request.TaskID = t.ID()
	if t.remainingCount <= 0 {
		t.isAlways = true
	}
	for _, s := range sinks {
		t.AddSink(s)
	}
	return t
}

func (f *Task) ID() string {
	return f.id
}

func (f *Task) AddSink(s sink.Sink) error {
	f.sinks = append(f.sinks, s)
	return nil
}

func (f *Task) Run() error {
	res, err := f.protocol.Handle(f.request, f.handlerOption)
	if err != nil {
		return err
	}

	fields, err := f.protocol.ContentMeta()
	if err != nil {
		return err
	}

	fields = append(fields, f.option.IncludeFields...)

	fc := NewFormatContent(&RawContent{Result: res, Request: f.request}, fields)
	for _, v := range f.sinks {
		go func(s sink.Sink, sc *string) {
			s.Run(sc)
		}(v, fc)
	}
	return nil
}

func (t *Task) RunLoop() error {
	ticker := time.NewTicker(t.interval)
	defer func() {
		ticker.Stop()
		log.Printf("task [%s] done.", t.ID())
	}()
	run := func() {
		if !t.isAlways {
			t.remainingCount--
		}
		if err := t.Run(); err != nil {
			log.Println(err)
		}
	}
	run()
	for t.remainingCount > 0 || t.isAlways {
		select {
		case <-t.ctx.Done():
			return nil
		case <-ticker.C:
			run()
		}
	}
	return nil
}

type TaskOption struct {
	IncludeFields []field.Field
	Count         int64
	Interval      time.Duration
}

type TaskContext struct {
	Request *core.Request
	Option  *TaskOption
	Sinks   []sink.Sink
}

func NewTaskContext(option *TaskOption, request *core.Request, sinkConfigs []SinkConfig) (*TaskContext, error) {
	sinks := []sink.Sink{}
	for _, sc := range sinkConfigs {
		sink, err := SinkOf(sc.Name, sc.Option)
		if err != nil {
			log.Println(err)
			continue
		}
		sinks = append(sinks, sink)
	}

	return &TaskContext{Request: request, Option: option, Sinks: sinks}, nil
}
