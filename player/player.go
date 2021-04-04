package player

import "net"

type Player struct {
	// TODO: How should we manage all of these connections? -> The player probably doesn't need it directly
	conn            *net.TCPConn
	State           State
	ProtocolVersion uint16
	Username        string
	Compression     CompressionState
}

type CompressionState struct {
	Enabled   bool
	Threshold uint64
}
