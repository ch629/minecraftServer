package packet

type (
	// Login State
	LoginStart struct {
		Name string
	}

	EncryptionResponse struct {
		SharedSecretLength int32  `pkt_type:"VarInt"`
		SharedSecret       []byte `pkt_len:"SharedSecretLength"`
		VerifyTokenLength  int32  `pkt_type:"VarInt"`
		VerifyToken        []byte `pkt_len:"VerifyTokenLength"`
	}

	LoginPluginResponse struct {
		MessageID  int32 `pkt_type:"VarInt"`
		Successful bool
		Data       []byte `pkt_opt:"Successful"`
	}
)
