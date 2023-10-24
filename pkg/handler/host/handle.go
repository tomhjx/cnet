package host

import (
	probing "github.com/prometheus-community/pro-bing"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/handler"
)

type Handle struct {
	Option *handler.Option
}

func (h Handle) Initialize(o *handler.Option) (handler.Handler, error) {
	return Handle{Option: o}, nil
}

func (h Handle) Do(hreq *core.CompletedRequest) (*core.Result, error) {
	res := core.NewResult()
	pinger, err := probing.NewPinger(hreq.ADDR)
	if err != nil {
		return res, err
	}
	pinger.Timeout = h.Option.TimeOut

	// pinger.OnRecv = func(pkt *probing.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	// }

	// pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	// }

	// pinger.OnFinish = func(stats *probing.Statistics) {
	// 	fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
	// 	fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
	// 		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	// 	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
	// 		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	// }

	if err := pinger.Run(); err != nil {
		return res, err
	}
	stats := pinger.Statistics()
	res.RunTime.Sent = stats.PacketsSent
	res.RunTime.Recv = stats.PacketsRecv
	res.RunTime.LossPct = stats.PacketLoss
	res.RunTime.StdDevTotalTime = stats.StdDevRtt
	res.RunTime.TotalTime = stats.AvgRtt
	return res, nil
}
