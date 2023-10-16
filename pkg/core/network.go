package core

type Network string

const (
	NetworkTCP        Network = "tcp"
	NetworkUnixSocket Network = "unixgram"
	NetworkUDP        Network = "udp"
)
