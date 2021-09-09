package game

import (
	"fmt"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DebugDisplayData struct {
	Enabled         bool
	TileDisplayMode int
}

const (
	DD_TILE_LIGHT                = iota
	DD_TILE_DISTANCE_FROM_PLAYER = iota
)

func handleTileDebugDisplay(tile *Tile, colour rl.Color) {
	switch state.DebugDisplay.TileDisplayMode {
	case DD_TILE_LIGHT:
		//! Tile light debug display
		rl.DrawText(fmt.Sprintf("%d", colour.A), tile.Pos.X, tile.Pos.Y, 12, rl.Red)
	case DD_TILE_DISTANCE_FROM_PLAYER:
		//! Tile distance debug display
		dist := getTileDistanceToPlayer(state.Player, tile)
		rl.DrawText(fmt.Sprintf("%.1f", dist), tile.Pos.X, tile.Pos.Y, 12, rl.Red)
	}

}

func drawDebugSettings() {
	if utils.DrawButton(rl.NewVector2(50.0, 110.0), "Tile light level") {
		state.DebugDisplay.TileDisplayMode = DD_TILE_LIGHT
	}

	if utils.DrawButton(rl.NewVector2(50.0, 230.0), "Tile distance from player") {
		state.DebugDisplay.TileDisplayMode = DD_TILE_DISTANCE_FROM_PLAYER
	}
}
