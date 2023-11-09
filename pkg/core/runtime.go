package core

import (
	"time"
)

type RunTime struct {
	NameLookUpTime      time.Duration
	ConnectTime         time.Duration
	TCPTime             time.Duration
	SSLTime             time.Duration
	AppConnectTime      time.Duration
	ServerProcessTime   time.Duration
	TTFB                time.Duration
	RedirectTime        time.Duration
	PreTransferTime     time.Duration
	StartTransferTime   time.Duration
	ContentTransferTime time.Duration
	TotalTime           time.Duration
	StdDevTotalTime     time.Duration
	ConnectedVia        string
	Sent                int
	Recv                int
	LossPct             float64
}
