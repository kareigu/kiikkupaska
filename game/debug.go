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

func handleTileDebugDisplay(tile *Tile) {
	switch state.DebugDisplay.TileDisplayMode {
	case DD_TILE_LIGHT:
		//! Tile light debug display
		rl.DrawText(fmt.Sprintf("%d", tile.LightLevel), tile.Pos.X, tile.Pos.Y, 12, rl.Red)
	case DD_TILE_DISTANCE_FROM_PLAYER:
		//! Tile distance debug display
		dist := tile.DistanceToPlayer()
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

	if utils.DrawButton(rl.NewVector2(50.0, 190.0), "Teleport to cursor") {
		state.Player.Pos = state.SelectionMode.Pos
	}

	if utils.DrawButton(rl.NewVector2(50.0, 220.0), "Spawn enemy on cursor") {
		nEnemy := CreateRandomEnemy(state.SelectionMode.Pos)
		state.Enemies = append(state.Enemies, nEnemy)
	}
}

func drawDebugInfo() {
	tileDebugInfo()
	enemiesDebugInfo()

	rl.DrawText(fmt.Sprintf("%.2f fps", rl.GetFPS()), 50, 20, 24, rl.RayWhite)
	rl.DrawText(fmt.Sprintf("%.4f ms", rl.GetFrameTime()*1000.0), 50, 50, 24, rl.RayWhite)
}

func tileDebugInfo() {
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

	data := fmt.Sprintf("Tile Pos: %v\nBlock: %v\nTile Type: %v", pos, tile.Block, tile.Type)

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
}

func enemiesDebugInfo() {
	enemyCount := 0
	var closestEnemy *Enemy

	for _, enemy := range state.Enemies {
		if enemy != nil {
			enemyCount++
			if closestEnemy == nil {
				closestEnemy = enemy
			} else {
				if closestEnemy.DistanceToPlayer() > enemy.DistanceToPlayer() {
					closestEnemy = enemy
				}
			}
		}
	}

	if closestEnemy != nil {
		pos := utils.IVector2{
			X: closestEnemy.Pos.X / TILE_SIZE,
			Y: closestEnemy.Pos.Y / TILE_SIZE,
		}

		p_pos := utils.IVector2{
			X: closestEnemy.LastKnownPlayerPos.X / TILE_SIZE,
			Y: closestEnemy.LastKnownPlayerPos.Y / TILE_SIZE,
		}
		data := fmt.Sprintf("Enemies in level: %v\nClosest Enemy: %.1f\nPos: %v\nLast player pos: %v", enemyCount, closestEnemy.DistanceToPlayer(), pos, p_pos)

		background := rl.NewRectangle(50.0, state.AppState.RES.ToVec2().Y-350.0, 250.0, 180.0)
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
	}
}
