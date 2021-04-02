package player

import "net"

type Player struct {
	conn            *net.TCPConn
	State           State
	ProtocolVersion int
	Username        string
	Compression     CompressionState
}

type CompressionState struct {
	Enabled   bool
	Threshold uint64
}
