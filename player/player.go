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

// CompressionState holds the defined compression options for a client connection
type CompressionState struct {
	// Enabled is whether compression is enabled for the client connection
	Enabled bool

	// Threshold is how many bytes each message can be before being compressed
	Threshold uint64
}
