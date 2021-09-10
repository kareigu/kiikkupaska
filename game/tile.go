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
			Type:  rendering.TILE_MISSING,
			Pos:   pos,
			Block: false,
		}
	}
}
