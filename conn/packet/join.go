package packet

import "github.com/Ixidi/flaming/conn"

type JoinGame struct {
	EntityId            conn.Int
	Hardcore            conn.Bool
	GameMode            conn.UByte
	PreviousGameMode    conn.UByte
	Worlds              conn.StringArray
	DimensionCodec      conn.NBT
	Dimension           conn.NBT
	WorldName           conn.String
	HashedSeed          conn.Long
	MaxPlayers          conn.VarInt
	ViewDistance        conn.VarInt
	ReducedDebugInfo    conn.Bool
	EnableReSpawnScreen conn.Bool
	Debug               conn.Bool
	Flat                conn.Bool
}

type DimensionCodecNBT struct {
	Dimensions DimensionsNBT `nbt:"minecraft:dimension_type"`
	Biomes     BiomesNBT     `nbt:"minecraft:worldgen/biome"`
}

type DimensionsNBT struct {
	Type  string                    `nbt:"type"`
	Value []IndexedDimensionTypeNBT `nbt:"value" nbt_type:"list"`
}

type IndexedDimensionTypeNBT struct {
	Name    string           `nbt:"name"`
	Id      int32            `nbt:"id"`
	Element DimensionTypeNBT `nbt:"element"`
}

type DimensionTypeNBT struct {
	AmbientLight       float32 `nbt:"ambient_light"`
	Infiniburn         string  `nbt:"infiniburn"`
	Natural            byte    `nbt:"natural"`
	HasCeiling         byte    `nbt:"has_ceiling"`
	HasSkylight        byte    `nbt:"has_skylight"`
	Ultrawarm          byte    `nbt:"ultrawarm"`
	HasRaids           byte    `nbt:"has_raids"`
	ReSpawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
	BedWorks           byte    `nbt:"bed_works"`
	PiglinSafe         byte    `nbt:"piglin_safe"`
	LogicalHeight      int32   `nbt:"logical_height"`
	CoordinateScale    int32   `nbt:"coordinate_scale"`
	Name               string  `nbt:"name"`
}

type BiomesNBT struct {
	Type  string            `nbt:"type"`
	Value []IndexedBiomeNBT `nbt:"value" nbt_type:"list"`
}

type IndexedBiomeNBT struct {
	Name    string   `nbt:"name"`
	Id      int32    `nbt:"id"`
	Element BiomeNBT `nbt:"element"`
}

type BiomeNBT struct {
	Depth         float32    `nbt:"depth"`
	Temperature   float32    `nbt:"temperature"`
	Scale         float32    `nbt:"scale"`
	Downfall      float32    `nbt:"downfall"`
	Category      string     `nbt:"category"`
	Precipitation string     `nbt:"precipitation"`
	Effects       EffectsNBT `nbt:"effects"`
}

type EffectsNBT struct {
	FogColor      int32 `nbt:"fog_color"`
	SkyColor      int32 `nbt:"sky_color"`
	WaterColor    int32 `nbt:"water_color"`
	WaterFogColor int32 `nbt:"water_fog_color"`
}
