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
		Count     int32       `pkt_type:"VarInt"`
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

	JoinGame struct {
		EntityID            int32
		IsHardcore          bool
		Gamemode            uint8
		PreviousGamemode    int8
		WorldCount          int32 `pkt_type:"VarInt"`
		WorldNames          []string
		DimensionCodec      DimensionCodecNBT `pkt_type:"nbt"`
		Dimension           DimensionTypeNBT  `pkt_type:"nbt"`
		WorldName           string
		HashedSeed          int64
		MaxPlayers          int32 `pkt_type:"VarInt"`
		ViewDistance        int32 `pkt_type:"VarInt"`
		ReducedDebugInfo    bool
		EnableRespawnScreen bool
		IsDebug             bool
		IsFlat              bool
	}

	DimensionCodecNBT struct {
		DimensionType struct {
			Type  string
			Value []struct {
				Name    string
				Id      int32
				Element DimensionTypeNBT
			}
		} `nbt:"minecraft:dimension_type"`
		Biome BiomeRegistryNBT `nbt:"minecraft:worldgen/biome"`
	}
	DimensionTypeNBT struct {
		PiglinSafe         bool `nbt:"piglin_safe"`
		Natural            bool
		AmbientLight       float32 `nbt:"ambient_light"`
		FixedTime          int64   `nbt:"fixed_time" nbt_opt:"true"`
		Infiniburn         string
		RespawnAnchorWorks bool `nbt:"respawn_anchor_works"`
		HasSkylight        bool `nbt:"has_skylight"`
		BedWorks           bool `nbt:"bed_works"`
		Effects            string
		HasRaids           bool    `nbt:"has_raids"`
		LogicalHeight      int32   `nbt:"logical_height"`
		CoordinateScale    float32 `nbt:"coordinate_scale"`
		Ultrawarm          bool
		HasCeiling         bool `nbt:"has_ceiling"`
	}
	BiomeRegistryNBT struct {
		Type  string
		Value struct {
			Name    string
			Id      int32
			Element struct {
				Precipitation string
				Depth         float32
				Temperature   float32
				Scale         float32
				Downfall      float32
				Category      string
				Effects       struct {
					SkyColor           int32  `nbt:"sky_color"`
					WaterFogColor      int32  `nbt:"water_fog_color"`
					FogColor           int32  `nbt:"fog_color"`
					WaterColor         int32  `nbt:"water_color"`
					FoliageColor       int32  `nbt:"foliage_color" nbt_opt:"true"`
					GrassColorModifier string `nbt:"grass_color_modifier" nbt_opt:"true"`
					Music              struct {
						ReplaceCurrentMusic bool `nbt:"replace_current_music"`
						Sound               string
						MaxDelay            int32 `nbt:"max_delay"`
						MinDelay            int32 `nbt:"min_delay"`
					}
					AmbientSound   string `nbt:"ambient_sound" nbt_opt:"true"`
					AdditionsSound struct {
						Sound      string
						TickChance float64 `nbt:"tick_chance"`
					} `nbt:"additions_sound" nbt_opt:"true"`
					MoodSound struct {
						Sound             string
						TickDelay         int32 `nbt:"tick_delay"`
						Offset            float64
						BlockSearchExtent int32 `nbt:"block_search_extent"`
					} `nbt:"mood_sound" nbt_opt:"true"`
					Particle struct {
						Probability float32
						Options     struct {
							Type string
						}
					} `nbt:"mood_sound" nbt_opt:"true"`
				}
			}
		}
	}
)
