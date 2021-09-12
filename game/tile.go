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
	Neighbours uint8
	LightLevel uint8
}

func (tile *Tile) Draw() {
	texture := rendering.GetTile(tile.Type)
	colour := rl.White

	if state.DebugDisplay.TileLightFx {
		colour.A = tile.LightLevel
	}

	if tile.Type == rendering.TILE_WALL_STONE {
		textures := &state.AppState.RenderAssets.TestTextures
		tile.UpdateNeighbours()

		rl.DrawTexture((*textures)[tile.Neighbours], tile.Pos.X, tile.Pos.Y, colour)
	} else {
		rl.DrawTexture(*texture, tile.Pos.X, tile.Pos.Y, colour)
	}
}

func (tile *Tile) UpdateNeighbours() {
	count := uint8(0)
	if nb, ok := GetMapTile(utils.NewIVector2(tile.Pos.X, tile.Pos.Y-TILE_SIZE)); ok && nb.Type != tile.Type {
		count += 1
	}
	if nb, ok := GetMapTile(utils.NewIVector2(tile.Pos.X+TILE_SIZE, tile.Pos.Y)); ok && nb.Type != tile.Type {
		count += 2
	}
	if nb, ok := GetMapTile(utils.NewIVector2(tile.Pos.X, tile.Pos.Y+TILE_SIZE)); ok && nb.Type != tile.Type {
		count += 4
	}
	if nb, ok := GetMapTile(utils.NewIVector2(tile.Pos.X-TILE_SIZE, tile.Pos.Y)); ok && nb.Type != tile.Type {
		count += 8
	}
	tile.Neighbours = count
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

func (tile *Tile) DistanceToEnemy(enemy *Enemy) float32 {
	tile_vec := tile.Pos.ToVec2()
	enemy_vec := enemy.Pos.ToVec2()
	distance := rl.Vector2Distance(tile_vec, enemy_vec) / float32(TILE_SIZE)

	return distance
}

func (tile *Tile) VisibleToPlayer(enemies *[]*Enemy) bool {
	distance := tile.DistanceToPlayer()
	tile.LightLevel = calculateLightLevel(distance, state.Player.Stats.Visibility)
	for _, enemy := range *enemies {
		if nlight := enemy.LightEmittedToTile(tile); nlight > tile.LightLevel {
			tile.LightLevel = nlight
		}
	}

	return tile.LightLevel > 0
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
