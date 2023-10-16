package flow

import "github.com/tomhjx/cnet/pkg/field"

func RegisterFields() {
	field.JobID.Inject(func(r any) any { return r.(*RawContent).Request.JobID })
	field.TaskID.Inject(func(r any) any { return r.(*RawContent).Request.TaskID })
	field.ClientID.Inject(func(r any) any { return r.(*RawContent).Request.ClientID })
	field.Tags.Inject(func(r any) any { return r.(*RawContent).Request.Tags })
	field.URL.Inject(func(r any) any { return r.(*RawContent).Request.RawURL })
	field.Method.Inject(func(r any) any { return r.(*RawContent).Request.Method })
	field.NameLookUpTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.NameLookUpTime.Milliseconds() })
	field.ConnectTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectTime.Milliseconds() })
	field.AppConnectTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.AppConnectTime.Milliseconds() })
	field.TCPTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TCPTime.Milliseconds() })
	field.SSLTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.SSLTime.Milliseconds() })
	field.TTFBTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TTFBTime.Milliseconds() })
	field.ServerProcessTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ServerProcessTime.Milliseconds() })
	field.PreTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.PreTransferTime.Milliseconds() })
	field.StartTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.StartTransferTime.Milliseconds() })
	field.ContentTransferTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ContentTransferTime.Milliseconds() })
	field.TotalTime.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.TotalTime.Milliseconds() })
	field.ConnectedVia.Inject(func(r any) any { return r.(*RawContent).Result.RunTime.ConnectedVia })
	field.Headers.Inject(func(r any) any { return r.(*RawContent).Result.Response.Headers })
	field.Body.Inject(func(r any) any { return r.(*RawContent).Result.Response.Body })
}
