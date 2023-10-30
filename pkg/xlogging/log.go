package xlogging

import "github.com/tomhjx/xlog"

func Configured() xlog.Verbose {
	return xlog.V(1)
}

func Changed() xlog.Verbose {
	return xlog.V(2)
}

func ChangedExtend() xlog.Verbose {
	return xlog.V(3)
}

func Debugged() xlog.Verbose {
	return xlog.V(4)
}

func Traced() xlog.Verbose {
	return xlog.V(5)
}
