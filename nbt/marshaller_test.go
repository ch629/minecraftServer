package nbt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalToNBTCases(t *testing.T) {
	type testCase struct {
		Name           string
		InputStruct    interface{}
		ExpectedOutput []byte
		Skip           bool
	}
	cases := []testCase{
		{
			Name: "Simple",
			InputStruct: struct {
				Root string
				Name string `nbt:"name"`
			}{
				Root: "hello world",
				Name: "Bananrama",
			},
			ExpectedOutput: []byte{
				0x0a,       // Compound
				0x00, 0x0b, // 11 Len
				0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, // name: hello world
				0x08,       // String
				0x00, 0x04, // 4 Len
				0x6e, 0x61, 0x6d, 0x65, // name: name
				0x00, 0x09, // 9 Len
				0x42, 0x61, 0x6e, 0x61, 0x6e, 0x72, 0x61, 0x6d, 0x61, // text: Bananrama
				0x00, // Tag End
			},
		},
		{
			Name: "Long List",
			InputStruct: struct {
				Root string
				List []int64 `nbt:"listTest (long)" nbt_type:"List"`
			}{
				List: []int64{11, 12, 13, 14, 15},
			},
			ExpectedOutput: []byte{
				0x0a,       // Compound
				0x00, 0x00, // 0 Len
				0x09,       // List
				0x00, 0x0f, // 15 Len
				0x6c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x20, 0x28, 0x6c, 0x6f, 0x6e, 0x67, 0x29, // listTest (long)
				0x04,                   // Long
				0x00, 0x00, 0x00, 0x05, // 5 Len
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, // 11
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, // 12
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0d, // 13
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e, // 14
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, // 15
				0x00, // Tag End
			},
		},
		{
			Name: "String List",
			InputStruct: struct {
				Strings []string
			}{
				Strings: []string{"hello", "world"},
			},
			ExpectedOutput: []byte{
				0x0a,       // Compound
				0x00, 0x00, // 0 Len
				0x09,       // List
				0x00, 0x07, // 7 Len
				0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, // strings
				0x08,                   // String
				0x00, 0x00, 0x00, 0x02, // 2 Len List
				0x00, 0x05, // 5 Len Str
				0x68, 0x65, 0x6c, 0x6c, 0x6f, // hello
				0x00, 0x05, // 5 Len Str
				0x77, 0x6f, 0x72, 0x6c, 0x64, // world
				0x00, // Tag End
			},
		},
		{
			Name: "Compound List",
			InputStruct: struct {
				Compounds []struct {
					Test uint8
				}
			}{
				Compounds: []struct{ Test uint8 }{
					{
						Test: 1,
					},
					{
						Test: 2,
					},
					{
						Test: 3,
					},
				},
			},
			ExpectedOutput: []byte{
				0x0a,       // Compound
				0x00, 0x00, // 0 Len
				0x09,       // List
				0x00, 0x09, // 9 Len
				0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x73, // compounds
				0x0a,                   // Compound
				0x00, 0x00, 0x00, 0x03, // 3 Len List
				0x01,       // Byte
				0x00, 0x04, // 4 Len
				0x74, 0x65, 0x73, 0x74, // test
				0x01,       // 1
				0x00,       // Tag End
				0x01,       // Byte
				0x00, 0x04, // 4 Len
				0x74, 0x65, 0x73, 0x74, // test
				0x02,       // 2
				0x00,       // Tag End
				0x01,       // Byte
				0x00, 0x04, // 4 Len
				0x74, 0x65, 0x73, 0x74, // test
				0x03, // 3
				0x00, // Tag End
				0x00, // Tag End
			},
		},
	}

	for _, test := range cases {
		if test.Skip {
			continue
		}
		bs, err := MarshalToNBT(test.InputStruct)
		assert.NoError(t, err, test.Name)
		assert.Equal(t, test.ExpectedOutput, bs, test.Name)
	}
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
	assert.Equal(t, []byte{0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x74, 0x79, 0x70, 0x65, 0x0, 0x0, 0xa, 0x0, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x0, 0x3, 0x0, 0x2, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x7, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0xa, 0x0, 0x0, 0x8, 0x0, 0xd, 0x70, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x0, 0x5, 0x0, 0x5, 0x64, 0x65, 0x70, 0x74, 0x68, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0xb, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x5, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x8, 0x64, 0x6f, 0x77, 0x6e, 0x66, 0x61, 0x6c, 0x6c, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x8, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x0, 0x0, 0xa, 0x0, 0x7, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x73, 0xa, 0x0, 0x0, 0x3, 0x0, 0x9, 0x73, 0x6b, 0x79, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xf, 0x77, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x66, 0x6f, 0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x9, 0x66, 0x6f, 0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xb, 0x77, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xd, 0x66, 0x6f, 0x6c, 0x69, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x14, 0x67, 0x72, 0x61, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x0, 0x0, 0xa, 0x0, 0x5, 0x6d, 0x75, 0x73, 0x69, 0x63, 0xa, 0x0, 0x0, 0x1, 0x0, 0x15, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x3, 0x0, 0x9, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x9, 0x6d, 0x69, 0x6e, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0xd, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0xa, 0x0, 0xf, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0xa, 0x0, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x6, 0x0, 0xb, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0xa, 0x6d, 0x6f, 0x6f, 0x64, 0x5f, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0xa, 0x0, 0x0, 0x8, 0x0, 0x5, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0x0, 0x0, 0x3, 0x0, 0xa, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0, 0x6, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x13, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0xa, 0x6d, 0x6f, 0x6f, 0x64, 0x5f, 0x73, 0x6f, 0x75, 0x6e, 0x64, 0xa, 0x0, 0x0, 0x5, 0x0, 0xb, 0x70, 0x72, 0x6f, 0x62, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x7, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0xa, 0x0, 0x0, 0x8, 0x0, 0x4, 0x74, 0x79, 0x70, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bs)
}
