package game

import (
	"math/rand"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	Type       int
	Pos        utils.IVector2
	Block      bool
	LightLevel uint8
}

func (tile *Tile) Draw() {
	texture := rendering.GetTile(tile.Type)
	colour := rl.White
	colour.A = tile.LightLevel
	rl.DrawTexture(*texture, tile.Pos.X, tile.Pos.Y, colour)
}

func (tile *Tile) Destroy() bool {
	if tile.Block {
		tile.Type = rendering.TILE_FLOOR_STONE
		tile.Block = false
		return true
	}
	return false
}

func (tile *Tile) DistanceToPlayer() float32 {
	tile_vec := tile.Pos.ToVec2()
	player_vec := state.Player.Pos.ToVec2()
	distance := rl.Vector2Distance(tile_vec, player_vec) / float32(TILE_SIZE)

	return distance
}

func (tile *Tile) VisibleToPlayer() (uint8, bool) {
	visrange := int32(state.Player.Stats.Visibility) * TILE_SIZE
	if tile.Pos.X > state.Player.Pos.X+(visrange) || tile.Pos.X < state.Player.Pos.X-(visrange) || tile.Pos.Y > state.Player.Pos.Y+(visrange) || tile.Pos.Y < state.Player.Pos.Y-(visrange) {
		return 0, false
	} else {
		distance := tile.DistanceToPlayer()
		return calculateLightLevel(distance, state.Player.Stats.Visibility), true
	}
}

func charToTile(c string, pos utils.IVector2) Tile {
	switch c {
	case "@":
		return Tile{
			Type:  rendering.TILE_WALL_STONE,
			Pos:   pos,
			Block: true,
		}
	case "_":
		ti := rendering.TILE_FLOOR_STONE
		if rand.Float32() < 0.01 {
			ti = rendering.TILE_FLOOR_STONE_BL
		}
		return Tile{
			Type:  ti,
			Pos:   pos,
			Block: false,
		}
	case "!":
		return Tile{
			Type:  rendering.TILE_WALL_MOSS,
			Pos:   pos,
			Block: true,
		}
	case "P":
		return Tile{
			Type:  rendering.TILE_FLOOR_SPAWN,
			Pos:   pos,
			Block: false,
		}
	case "-":
		return Tile{
			Type:  rendering.TILE_FLOOR_OBS,
			Pos:   pos,
			Block: false,
		}
	default:
		return Tile{
			Type:  rendering.TILE_FLOOR_STONE,
			Pos:   pos,
			Block: false,
		}
	}
}
