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
	TTFBTime            time.Duration
	RedirectTime        time.Duration
	PreTransferTime     time.Duration
	StartTransferTime   time.Duration
	ContentTransferTime time.Duration
	TotalTime           time.Duration
	ConnectedVia        string
}
