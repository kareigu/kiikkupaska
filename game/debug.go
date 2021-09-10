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
	DD_TILE_NO_DISPLAY           = iota
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
	if utils.DrawButton(rl.NewVector2(50.0, 100.0), "No display") {
		state.DebugDisplay.TileDisplayMode = DD_TILE_NO_DISPLAY
	}

	if utils.DrawButton(rl.NewVector2(50.0, 130.0), "Tile light level") {
		state.DebugDisplay.TileDisplayMode = DD_TILE_LIGHT
	}

	if utils.DrawButton(rl.NewVector2(50.0, 160.0), "Tile distance from player") {
		state.DebugDisplay.TileDisplayMode = DD_TILE_DISTANCE_FROM_PLAYER
	}
}

func drawDebugInfo() {
	var tile *Tile
	var sourcePos utils.IVector2
	if state.SelectionMode.Using {
		sourcePos = state.SelectionMode.Pos
	} else {
		sourcePos = state.Player.Pos
	}

	if t, ok := getTile(sourcePos); ok {
		tile = t
	}

	background := rl.NewRectangle(50.0, 280.0, 250.0, 180.0)

	pos := utils.IVector2{
		X: tile.Pos.X / TILE_SIZE,
		Y: tile.Pos.Y / TILE_SIZE,
	}

	data := fmt.Sprintf("Tile Pos: %v\nBlock: %v\nTexture ID: %v", pos, tile.Block, tile.Texture.ID)

	rl.DrawRectangleRec(background, rl.DarkGray)

	rl.DrawTextRec(
		state.AppState.RenderAssets.SecondaryFont,
		data,
		background,
		24.0,
		1.0,
		true,
		rl.White,
	)

	rl.DrawText(fmt.Sprintf("%.2f fps", rl.GetFPS()), 50, 20, 24, rl.RayWhite)
	rl.DrawText(fmt.Sprintf("%.4f ms", rl.GetFrameTime()*1000.0), 50, 50, 24, rl.RayWhite)
}
