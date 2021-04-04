package main

import (
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
	LoginData struct {
		Payload string `json:"payload"`
	}

	HandshakeData struct {
		ProtocolVersion int32  `json:"protocol_version" pkt_type:"VarInt"`
		ServerAddress   string `json:"server_address"`
		ServerPort      uint16 `json:"server_port"`
		NextState       int32  `json:"next_state" pkt_type:"VarInt"`
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
	var handshake HandshakeData
	err := packet.Unmarshal(pkt, &handshake)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal HandshakeData")
	}

	return &handshake, nil
}

func ReadLoginData(pkt packet.Packet) (*LoginData, error) {
	var loginData LoginData
	if err := packet.Unmarshal(pkt, &loginData); err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal LoginData")
	}

	return &loginData, nil
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
				fmt.Println("OPENING CONNECTION")
				defer func() {
					fmt.Println("CLOSING CONNECTION")
				}()
				state := int32(0)
				for {
					pkt, err := packet.MakeUncompressedPacket(conn)
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
						state = h.NextState
						fmt.Println(h)
					} else if state == 2 {
						loginData, err := ReadLoginData(pkt)
						if IsConnectionClosedErr(err) {
							break
						}
						p(err)
						fmt.Println(loginData)

						// Login success testing -> We get 'Joining world' from this
						loginSuccess := &packet.LoginSuccess{
							UUID:     uuid.MustParse("e52d49e2f2244a7380cfcacf6aecbcae"),
							Username: loginData.Payload,
						}
						newPkt, err := packet.MakePacketWithData(0x02, loginSuccess)
						p(err)
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
