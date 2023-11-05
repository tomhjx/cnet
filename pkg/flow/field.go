package flow

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/metric"
)

func RegisterFields() {

	field.ID.InitValueOf(func(r any) any { return r.(*RawContent).Request.ID })
	field.JobID.InitValueOf(func(r any) any { return r.(*RawContent).Request.JobID })
	field.TaskID.InitValueOf(func(r any) any { return r.(*RawContent).Request.TaskID })
	field.ClientID.InitValueOf(func(r any) any { return r.(*RawContent).Request.ClientID })
	field.Tags.InitValueOf(func(r any) any { return r.(*RawContent).Request.Tags })
	field.ADDR.InitValueOf(func(r any) any { return r.(*RawContent).Request.ADDR })
	field.Method.InitValueOf(func(r any) any { return r.(*RawContent).Request.Method })
	field.NameLookUpTime.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.NameLookUpTime.Seconds() })
	field.ConnectTime.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectTime.Seconds() })
	field.AppConnectTime.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.AppConnectTime.Seconds() })

	field.TCPTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.TCPTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.SSLTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.SSLTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.TTFBTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.TTFBTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.ServerProcessTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.ServerProcessTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.PreTransferTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.PreTransferTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.StartTransferTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.StartTransferTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.ContentTransferTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.ContentTransferTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.TotalTime.Init(&field.FieldHandle{
		ValueOf: func(r any) any { return r.(*RawContent).Result.RunTime.TotalTime.Seconds() },
		Metric:  defaultMetricDecorate(),
	})

	field.ConnectedVia.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectedVia })
	field.Headers.InitValueOf(func(r any) any { return r.(*RawContent).Result.Response.Headers })
	field.Body.InitValueOf(func(r any) any { return r.(*RawContent).Result.Response.Body })
	field.Status.InitValueOf(func(r any) any { return r.(*RawContent).Result.Response.Status })
	field.StatusCode.InitValueOf(func(r any) any { return r.(*RawContent).Result.Response.StatusCode })
	field.Sent.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.Sent })
	field.Recv.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.Recv })
	field.LossPct.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.LossPct })
	field.StdDevTotalTime.InitValueOf(func(r any) any { return r.(*RawContent).Result.RunTime.StdDevTotalTime.Seconds() })
}

func fillMetricLabels(labels prometheus.Labels, r any) prometheus.Labels {
	raw := r.(*RawContent)
	labels["code"] = fmt.Sprint(raw.Result.Response.StatusCode)
	labels["method"] = raw.Request.Method
	return labels
}

func defaultMetricDecorate() *field.Metric {
	return &field.Metric{
		LabelsOf: fillMetricLabels,
		Classes:  []metric.Metric{metric.DistributeMetric},
	}
}
