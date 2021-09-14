package game

import (
	"fmt"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleControls() {
	if !state.Player.Turn.Done {
		if state.UIState.SelectionMode.Using {
			moveSelectionCursor(&state.UIState.SelectionMode)
		} else {
			state.Player.Move()
		}
		state.tempTimeSinceTurn = 0.0
	} else {
		state.UIState.SelectionMode.Using = false

		enemyTurnsComplete := false
		for _, enemy := range state.Enemies {
			enemyTurnsComplete = enemy.Turn.Done == true
		}
		if enemyTurnsComplete {
			fmt.Printf("Enemy turns processed in %.3f ms\n", state.tempTimeSinceTurn)
			state.Player.StartTurn()
		} else {
			state.tempTimeSinceTurn += rl.GetFrameTime()
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		state.UIState.SelectionMode.Pos = state.Player.Pos
		state.UIState.SelectionMode.Using = !state.UIState.SelectionMode.Using
	}

	if rl.IsKeyPressed(rl.KeyM) || rl.IsKeyPressed(rl.KeyEscape) {
		state.AppState.View = utils.PAUSED
	}

	zoomMult := float32(rl.GetMouseWheelMove()) * 0.1
	if state.Camera.Zoom+zoomMult > 0.2 {
		state.Camera.Zoom += zoomMult
	}

	if state.UIState.SelectionMode.Using {
		if state.Player.Turn.Actions > 0 {
			if rl.IsKeyPressed(rl.KeyB) {
				if tile, ok := GetMapTile(state.UIState.SelectionMode.Pos); ok {
					if tile.Destroy() {
						state.Player.Turn.Actions--
					}
				}
			}
			if rl.IsKeyPressed(rl.KeyV) {
				for _, enemy := range state.Enemies {
					if enemy.Pos == state.UIState.SelectionMode.Pos {
						state.Player.Attack(enemy)
					}
				}
			}
		}
	}

	if utils.DebugMode {
		if rl.IsKeyPressed(rl.KeyF1) {
			state.UIState.DebugDisplay.Enabled = !state.UIState.DebugDisplay.Enabled
		}

		if rl.IsKeyPressed(rl.KeyI) {
			state.Player.Stats.Visibility++
		}

		if rl.IsKeyPressed(rl.KeyK) {
			state.Player.Stats.Visibility--
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		state.Player.EndTurn()
	}

	state.Camera.Target.X = float32(state.Player.Pos.X)
	state.Camera.Target.Y = float32(state.Player.Pos.Y)
}

func moveSelectionCursor(selection *SelectionMode) {
	s_x := selection.Pos.X
	s_y := selection.Pos.Y

	if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
		s_x -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
		s_x += TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
		s_y -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		s_y += TILE_SIZE
	}

	selection.Pos.X = s_x
	selection.Pos.Y = s_y
}
