package field

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tomhjx/cnet/pkg/metric"
)

type Field string

var fields = map[Field]*FieldHandle{}

func (f Field) ValueOf(r any) (any, error) {
	h, err := f.handle()
	if err != nil {
		return r, err
	}
	return h.ValueOf(r), nil
}

func (f Field) CalcMetric(labels prometheus.Labels, r any) error {
	h, err := f.handle()
	if err != nil {
		return err
	}
	if h.Metric == nil {
		return nil
	}
	v, err := f.ValueOf(r)
	if err != nil {
		return err
	}
	labels["field"] = string(f)

	if h.Metric.LabelsOf != nil {
		labels = h.Metric.LabelsOf(labels, r)
	}
	rv := v.(float64)
	for _, m := range h.Metric.Classes {
		m.Calc(labels, rv)
	}
	return nil
}

func (f Field) ResetMetric() error {
	h, err := f.handle()
	if err != nil {
		return err
	}
	for _, m := range h.Metric.Classes {
		m.Reset()
	}
	return nil
}

func (f Field) IsEnableMetric() bool {
	h, err := f.handle()
	if err != nil {
		return false
	}
	if h.Metric == nil {
		return false
	}
	return true
}

func (f Field) String() string {
	return string(f)
}

func (f Field) handle() (*FieldHandle, error) {
	if h, ok := fields[f]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("not support field [%s]", f)
}

func (f Field) Init(h *FieldHandle) *FieldHandle {
	fields[f] = h
	return h
}

func (f Field) InitValueOf(valueOf func(r any) any) *FieldHandle {
	h := &FieldHandle{
		ValueOf: valueOf,
	}
	f.Init(h)
	return h
}

type FieldValueHandle func(any) any
type Metric struct {
	LabelsOf func(prometheus.Labels, any) prometheus.Labels
	Classes  []metric.Metric
}

type FieldHandle struct {
	ValueOf func(r any) any
	Metric  *Metric
}
