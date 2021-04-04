package packet

import "github.com/google/uuid"

type (
	// Login State
	Disconnect struct {
		Reason string
	}

	EncryptionRequest struct {
		ServerID          string
		PublicKeyLength   int32  `pkt_type:"VarInt"`
		PublicKey         []byte `pkt_len:"PublicKeyLength"`
		VerifyTokenLength int32  `pkt_type:"VarInt"`
		VerifyToken       []byte `pkt_len:"VerifyTokenLength"`
	}

	LoginSuccess struct {
		// TODO: Check if this UUID works
		UUID     uuid.UUID
		Username string
	}

	SetCompression struct {
		Threshold int32 `pkt_type:"VarInt"`
	}

	LoginPluginRequest struct {
		MessageID int32 `pkt_type:"VarInt"`
		Channel   string
		Data      []byte
	}

	// Play State
	SpawnEntity struct {
		EntityID   int32 `pkt_type:"VarInt"`
		ObjectUUID uuid.UUID
		Type       int32 `pkt_type:"VarInt"`
		X          float64
		Y          float64
		Z          float64
		// TODO: Check Angle works
		Pitch     Angle
		Yaw       Angle
		Data      int32
		VelocityX int16
		VelocityY int16
		VelocityZ int16
	}

	SpawnExperienceOrb struct {
		EntityID int32 `pkt_type:"VarInt"`
		X        float64
		Y        float64
		Z        float64
		Count    int16
	}

	SpawnLivingEntity struct {
		EntityID   int32 `pkt_type:"VarInt"`
		EntityUUID uuid.UUID
		X          float64
		Y          float64
		Z          float64
		Yaw        Angle
		Pitch      Angle
		HeadPitch  Angle
		VelocityX  int16
		VelocityY  int16
		VelocityZ  int16
	}

	SpawnPainting struct {
		EntityID   int32 `pkt_type:"VarInt"`
		EntityUUID uuid.UUID
		Motive     int32 `pkt_type:"VarInt"`
		Location   Position
		// TODO: Direction wrapper
		Direction byte
	}

	SpawnPlayer struct {
		EntityID   int32 `pkt_type:"VarInt"`
		PlayerUUID uuid.UUID
		X          float64
		Y          float64
		Z          float64
		Yaw        Angle
		Pitch      Angle
	}

	EntityAnimation struct {
		EntityID int32 `pkt_type:"VarInt"`
		// TODO: Animation wrapper
		Animation uint8
	}

	Statistics struct {
		Count int32 `pkt_type:"VarInt"`
		// TODO: Array of structs need to be implemented
		Statistic []Statistic `pkt_len:"Count"`
	}

	Statistic struct {
		CategoryID  int32 `pkt_type:"VarInt"`
		StatisticID int32 `pkt_type:"VarInt"`
		Value       int32 `pkt_type:"VarInt"`
	}

	AcknowledgePlayerDigging struct {
		Location   Position
		Block      int32 `pkt_type:"VarInt"`
		Status     int32 `pkt_type:"VarInt"`
		Successful bool
	}

	BlockBreakAnimation struct {
		EntityID     int32 `pkt_type:"VarInt"`
		Location     Position
		DestroyStage byte
	}
)
