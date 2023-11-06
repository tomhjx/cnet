package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomhjx/cnet/pkg/xlogging"
	"github.com/tomhjx/xlog"
)

var (
	metrics = map[Metric]*MetricHandle{}
	enable  = true
)

func Disable() {
	enable = false
}
func Enable() {
	enable = true
}

func IsEnabled() bool {
	return enable
}

type MetricHandle struct {
	Collector  prometheus.Collector
	LabelNames []string
	Calc       func(labels prometheus.Labels, v float64)
	Reset      func()
}

type Metric string

func (m Metric) Init(mc *MetricHandle) {
	prometheus.MustRegister(mc.Collector)
	metrics[m] = mc
}

func (m Metric) Calc(labels prometheus.Labels, v float64) error {
	mc, err := m.handle()
	if err != nil {
		return err
	}
	for _, v := range mc.LabelNames {
		if _, ok := labels[v]; !ok {
			labels[v] = ""
		}
	}
	mc.Calc(labels, v)
	return nil
}

func (m Metric) Reset() error {
	mc, err := m.handle()
	if err != nil {
		return err
	}
	mc.Reset()
	return nil
}

func (mc Metric) handle() (*MetricHandle, error) {
	if m, ok := metrics[mc]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("not support metric class [%s]", mc)
}

const (
	TotalMetric      Metric = "cnet_flow_total"
	ValueMetric      Metric = "cnet_flow_value"
	DistributeMetric Metric = "cnet_flow_distribute"
)

func RegisterMetrics() {

	labelNames := []string{"addr", "field", "method", "code"}

	totalMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: string(TotalMetric),
	}, labelNames)

	TotalMetric.Init(&MetricHandle{
		Collector:  totalMetric,
		LabelNames: labelNames,
		Calc:       func(labels prometheus.Labels, v float64) { totalMetric.With(labels).Add(v) },
		Reset:      func() { totalMetric.Reset() },
	})

	valueMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: string(ValueMetric),
	}, labelNames)
	ValueMetric.Init(&MetricHandle{
		Collector:  valueMetric,
		LabelNames: labelNames,
		Calc:       func(labels prometheus.Labels, v float64) { valueMetric.With(labels).Set(v) },
		Reset:      func() { valueMetric.Reset() },
	})

	distributeMetric := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       string(DistributeMetric),
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labelNames)
	DistributeMetric.Init(&MetricHandle{
		Collector:  distributeMetric,
		LabelNames: labelNames,
		Calc:       func(labels prometheus.Labels, v float64) { distributeMetric.With(labels).Observe(v) },
		Reset:      func() { distributeMetric.Reset() },
	})

}

func StartServer(port int, ctx context.Context) {
	if !IsEnabled() {
		return
	}

	RegisterMetrics()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: mux,
	}
	go func() {
		xlogging.Changed().Infof("metrics server(:%d) started.", port)
		if err := srv.ListenAndServe(); err != nil {
			xlog.Error(err, "metrics server start fail.")
		}
	}()
}
