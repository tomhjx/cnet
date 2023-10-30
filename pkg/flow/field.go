package flow

import "github.com/tomhjx/cnet/pkg/field"

func RegisterFields() {
	field.ID.Inject(func(r any) any { return r.(*RawContent).Request.ID })
	field.JobID.Inject(func(r any) any { return r.(*RawContent).Request.JobID })
	field.TaskID.Inject(func(r any) any { return r.(*RawContent).Request.TaskID })
	field.ClientID.Inject(func(r any) any { return r.(*RawContent).Request.ClientID })
	field.Tags.Inject(func(r any) any { return r.(*RawContent).Request.Tags })
	field.ADDR.Inject(func(r any) any { return r.(*RawContent).Request.ADDR })
	field.Method.Inject(func(r any) any { return r.(*RawContent).Request.Method })
	field.NameLookUpTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.NameLookUpTime.Seconds() })
	field.ConnectTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectTime.Seconds() })
	field.AppConnectTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.AppConnectTime.Seconds() })
	field.TCPTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TCPTime.Seconds() })
	field.SSLTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.SSLTime.Seconds() })
	field.TTFBTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TTFBTime.Seconds() })
	field.ServerProcessTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ServerProcessTime.Seconds() })
	field.PreTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.PreTransferTime.Seconds() })
	field.StartTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.StartTransferTime.Seconds() })
	field.ContentTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ContentTransferTime.Seconds() })
	field.TotalTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TotalTime.Seconds() })
	field.ConnectedVia.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectedVia })
	field.Headers.Inject(func(r any) any { return r.(*RawContent).Result.Response.Headers })
	field.Body.Inject(func(r any) any { return r.(*RawContent).Result.Response.Body })
	field.Status.Inject(func(r any) any { return r.(*RawContent).Result.Response.Status })
	field.StatusCode.Inject(func(r any) any { return r.(*RawContent).Result.Response.StatusCode })
	field.Sent.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.Sent })
	field.Recv.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.Recv })
	field.LossPct.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.LossPct })
	field.StdDevTotalTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.StdDevTotalTime.Seconds() })
}
