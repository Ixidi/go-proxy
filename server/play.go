package server

import (
	"github.com/Ixidi/flaming/conn"
	"github.com/Ixidi/flaming/conn/packet"
)

func (s *Server) spawnPlayer(player *Player) error {
	dimensionTypeNbt := packet.DimensionTypeNBT{
		AmbientLight:       0.0,
		Infiniburn:         "",
		Natural:            1,
		HasCeiling:         1,
		HasSkylight:        1,
		Ultrawarm:          0,
		HasRaids:           1,
		ReSpawnAnchorWorks: 0,
		BedWorks:           1,
		PiglinSafe:         0,
		LogicalHeight:      256,
		CoordinateScale:    1,
		Name:               "minecraft:overworld",
	}

	indexedDimensionNbt := packet.IndexedDimensionTypeNBT{
		Name:    "minecraft:overworld",
		Id:      0,
		Element: dimensionTypeNbt,
	}

	biomeNbt := packet.BiomeNBT{
		Depth:         0.125,
		Temperature:   0.8,
		Scale:         0.05,
		Downfall:      0.4,
		Category:      "none",
		Precipitation: "none",
		Effects: packet.EffectsNBT{
			FogColor:      0xC0D8FF,
			SkyColor:      0x78A7FF,
			WaterColor:    0x3F76E4,
			WaterFogColor: 0x50533,
		},
	}

	indexedBiomeNbt := packet.IndexedBiomeNBT{
		Name:    "minecraft:plains",
		Id:      0,
		Element: biomeNbt,
	}

	codecNbt := packet.DimensionCodecNBT{
		Dimensions: packet.DimensionsNBT{
			Type:  "minecraft:dimension_type",
			Value: []packet.IndexedDimensionTypeNBT{indexedDimensionNbt},
		},
		Biomes: packet.BiomesNBT{
			Type:  "minecraft:worldgen/biome",
			Value: []packet.IndexedBiomeNBT{indexedBiomeNbt},
		},
	}

	p, err := conn.Pack(packet.JoinGameId, &packet.JoinGame{
		EntityId:            0,
		Hardcore:            false,
		GameMode:            0,
		PreviousGameMode:    0,
		Worlds:              []conn.String{"world"},
		DimensionCodec:      conn.NBT{V: codecNbt},
		Dimension:           conn.NBT{V: dimensionTypeNbt},
		WorldName:           "world",
		HashedSeed:          0,
		MaxPlayers:          1,
		ViewDistance:        2,
		ReducedDebugInfo:    false,
		EnableReSpawnScreen: true,
		Debug:               false,
		Flat:                false,
	})
	if err != nil {
		return err
	}
	player.outgoingPackets <- p

	p, err = conn.Pack(packet.SpawnPositionId, &packet.SpawnPosition{
		Position: conn.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	})
	if err != nil {
		return err
	}

	player.outgoingPackets <- p

	p, err = conn.Pack(packet.PlayerPositionAndLookId, &packet.PlayerPositionAndLook{
		X:          0,
		Y:          0,
		Z:          0,
		Yaw:        0,
		Pitch:      0,
		Flags:      0,
		TeleportId: 0,
	})
	if err != nil {
		return err
	}
	player.outgoingPackets <- p

	return nil
}
