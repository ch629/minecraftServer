package nbt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalToNBT(t *testing.T) {
	type Compound struct {
		Root string
		Name string `nbt:"name"`
	}

	bs, err := MarshalToNBT(&Compound{
		Root: "hello world",
		Name: "Bananrama",
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte{
		0x0a,       // Compound
		0x00, 0x0b, // 11 Len
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // name: hello world
		0x08,       // String
		0x00, 0x04, // 4 Len
		0x6e, 0x61, 0x6d, 0x65, // name: name
		0x00, 0x09, // 9 Len
		0x42, 0x61, 0x6e, 0x61, 0x6e, 0x72, 0x61, 0x6d, 0x61, // text: Bananrama
		0x00, // Tag End
	}, bs)
}

func TestMarshalToNBT2(t *testing.T) {
	type BiomeRegistryNBT struct {
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
	n := BiomeRegistryNBT{}
	bs, err := MarshalToNBT(&n)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x74, 0x79, 0x70, 0x65, 0x0, 0x0, 0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x0, 0x3, 0x0, 0x2, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x8, 0x0, 0xd, 0x70, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x0, 0x5, 0x0, 0x5, 0x64, 0x65, 0x70, 0x74, 0x68, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0xb, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x5, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x8, 0x64, 0x6f, 0x77, 0x6e, 0x66, 0x61, 0x6c, 0x6c, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x8, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x0, 0x0, 0xa, 0x0, 0x0, 0x3, 0x0, 0x9, 0x73, 0x6b, 0x79, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xf, 0x77, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x66, 0x6f, 0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x9, 0x66, 0x6f, 0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xb, 0x77, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xd, 0x66, 0x6f, 0x6c, 0x69, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x14, 0x67, 0x72, 0x61, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x0, 0x0, 0xa, 0x0, 0x0, 0x1, 0x0, 0x15, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x3, 0x0, 0x9, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x9, 0x6d, 0x69, 0x6e, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0xd, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0xa, 0x0, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x6, 0x0, 0xb, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x3, 0x0, 0xa, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0, 0x6, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x13, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x5, 0x0, 0xb, 0x70, 0x72, 0x6f, 0x62, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x74, 0x79, 0x70, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bs)
}