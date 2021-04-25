package connection

import (
	"context"
	"github.com/ch629/minecraftServer/packet"
	"github.com/rotisserie/eris"
	"io"
	"net"
)

// TODO: Do we need to hold the compression state here, so it'll automatically wrap our reader?
type (
	ChanneledConnection interface {
		// Input returns a channel of the data coming from a client connection
		Input() <-chan packet.Packet

		// Errors returns a channel of errors produced when reading data from a client connection
		Errors() <-chan error

		// Send sends packet.Packet data to the client connection synchronously
		Send(pkt packet.Packet) error

		// Close closes the client connection, channels & stops goroutines
		Close()
	}

	channeledConnection struct {
		ctx       context.Context
		ctxCancel context.CancelFunc
		conn      *net.TCPConn
		input     chan packet.Packet
		errors    chan error
	}
)

var (
	// ErrConnClosed is an error to notify when the connection has been closed
	// This error is sent to the Errors channel when the connection is initially closed
	// and returned from Send
	ErrConnClosed = eris.New("connection was closed")
)

func (con *channeledConnection) Input() <-chan packet.Packet {
	return con.input
}

func (con *channeledConnection) Errors() <-chan error {
	return con.errors
}

func (con *channeledConnection) Close() {
	select {
	case <-con.ctx.Done():
		return
	default:
	}
	con.ctxCancel()
	con.error(con.conn.Close())
	con.error(ErrConnClosed)
	close(con.errors)
	close(con.input)
}

// Send sends a packet.Packet to the client connection, returning ErrConnClosed if the connection was closed previously
func (con *channeledConnection) Send(pkt packet.Packet) error {
	select {
	case <-con.ctx.Done():
		return ErrConnClosed
	default:
	}
	_, err := packet.WriteTo(pkt, con.conn)
	return err
}

func (con *channeledConnection) inputLoop() {
	for {
		select {
		case <-con.ctx.Done():
			return
		default:
		}

		pkt, err := packet.MakeUncompressedPacket(con.conn)
		if err != nil {
			if eris.Is(err, io.EOF) {
				con.Close()
				return
			}
			con.error(err)
			continue
		}

		// TODO: Make sure we never get to this point if the channel is closed
		con.input <- pkt
	}
}

// error writes an error to the channel without blocking if err is not nil
func (con *channeledConnection) error(err error) {
	if err != nil {
		select {
		case <-con.ctx.Done():
		case con.errors <- err:
		default:
		}
	}
}

func MakeChanneledConnection(ctx context.Context, conn *net.TCPConn) ChanneledConnection {
	con := &channeledConnection{
		conn:   conn,
		input:  make(chan packet.Packet),
		errors: make(chan error),
	}
	con.ctx, con.ctxCancel = context.WithCancel(ctx)
	go con.inputLoop()
	return con
}
