package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"io"
	"minecraftServer/packet"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type (
	HandshakeData struct {
		ProtocolVersion packet.VarInt `json:"protocol_version"`
		ServerAddress   packet.String `json:"server_address"`
		ServerPort      packet.UShort `json:"server_port"`
		NextState       State         `json:"next_state"`
	}

	LoginData struct {
		Payload packet.String `json:"payload"`
	}

	State byte
)

const (
	Handshaking State = iota
	Status
	Login
	Play
)

func StateFromVarInt(varInt packet.VarInt) State {
	return []State{Handshaking, Status, Login, Play}[varInt]
}

func (s State) String() string {
	return []string{"Handshaking", "Status", "Login", "Play"}[s]
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (d HandshakeData) String() string {
	b, _ := json.Marshal(d)
	return fmt.Sprintf("Handshake %v", string(b))
}

func (d LoginData) String() string {
	b, _ := json.Marshal(d)
	return fmt.Sprintf("LoginData %v", string(b))
}

func ReadHandshakeData(pkt packet.Packet) (*HandshakeData, error) {
	reader, err := pkt.DataReader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var version packet.VarInt
	var addr packet.String
	var port packet.UShort
	var nextState packet.VarInt

	if err = packet.ReadFields(reader, &version, &addr, &port, &nextState); err != nil {
		return nil, eris.Wrap(err, "failed to read handshake data")
	}

	return &HandshakeData{
		ProtocolVersion: version,
		ServerAddress:   addr,
		ServerPort:      port,
		NextState:       StateFromVarInt(nextState),
	}, nil
}

func ReadLoginData(pkt packet.Packet) (*LoginData, error) {
	reader, err := pkt.DataReader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var payload packet.String
	if err := packet.ReadFields(reader, &payload); err != nil {
		return nil, eris.Wrap(err, "failed to read login data")
	}

	return &LoginData{Payload: payload}, nil
}

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:25565")
	p(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	p(err)
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			p(err)
			go func(conn *net.TCPConn) {
				defer conn.Close()
				defer func() {
					fmt.Println("CLOSING CONNECTION")
				}()
				tee := io.TeeReader(conn, os.Stdout)
				state := int32(0)
				for {
					pkt, err := packet.MakeUncompressedPacket(tee)
					if IsConnectionClosedErr(err) {
						break
					}
					p(err)
					if state == 0 && pkt.ID() == 0 {
						h, err := ReadHandshakeData(pkt)
						if IsConnectionClosedErr(err) {
							break
						}
						p(err)
						state = int32(h.NextState)
						fmt.Println(h)
					} else if state == 2 {
						loginData, err := ReadLoginData(pkt)
						if IsConnectionClosedErr(err) {
							break
						}
						p(err)
						fmt.Println(loginData)

						// Login success testing -> We get 'Joining world' from this
						dataBuf := bytes.NewBuffer(nil)
						uuid, _ := uuid.NewRandom()
						err = packet.WriteFields(dataBuf, packet.UUID(uuid), loginData.Payload)
						p(err)
						newPkt := packet.MakePacket(packet.VarInt(0x02), dataBuf)
						_, err = packet.WriteTo(newPkt, conn)
						p(err)
					}
				}
			}(conn)
		}
	}()
	<-sigs
}

func IsConnectionClosedErr(err error) bool {
	return err == io.EOF || eris.Unwrap(err) == io.EOF
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
