package player

import (
	"encoding/json"
	"github.com/ch629/minecraftServer/packet"
)

type (
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
