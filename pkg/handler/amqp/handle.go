package amqp

import (
	"context"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/handler"
	"github.com/tomhjx/xlog"
)

type Handle struct {
	Option *handler.Option
}

func (h Handle) Initialize(o *handler.Option) (handler.Handler, error) {
	return Handle{Option: o}, nil
}

func (h Handle) Do(hreq *core.CompletedRequest) (*core.Result, error) {

	res := core.NewResult()

	t0 := time.Now()
	conn, err := amqp.Dial(hreq.ADDR)
	if err != nil {
		return nil, err
	}
	res.RunTime.ConnectTime = time.Since(t0)
	t1 := time.Now()
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgID := uuid.NewString()
	cMsgs, err := ch.Consume(
		hreq.Queue, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}
	res.RunTime.PreTransferTime = time.Since(t1)
	t2 := time.Now()
	err = ch.PublishWithContext(ctx, hreq.QueueExchange, hreq.QueueRoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("aaa"),
		Timestamp:   time.Now(),
		MessageId:   msgID,
	})
	if err != nil {
		return nil, err
	}
	t3 := time.Now()
	res.RunTime.StartTransferTime = time.Since(t2)
	for d := range cMsgs {
		if d.MessageId != msgID {
			continue
		}
		xlog.Info(d.Timestamp, time.Now())
		res.RunTime.TTFBTime = time.Since(t3)
		break
	}
	res.RunTime.TotalTime = time.Since(t0)
	return res, nil
}
