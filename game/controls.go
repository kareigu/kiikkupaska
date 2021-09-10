package game

import (
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleControls() {
	if !state.Player.Turn.Done {
		if state.SelectionMode.Using {
			moveSelectionCursor(&state.SelectionMode)
		} else {
			movePlayer(state.Player, &state.Map)
		}
		state.tempTimeSinceTurn = 0.0
	} else {
		state.SelectionMode.Using = false
		if state.tempTimeSinceTurn > 2.0 {
			state.Player.StartTurn()
		} else {
			state.tempTimeSinceTurn += rl.GetFrameTime()
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		state.SelectionMode.Pos = state.Player.Pos
		state.SelectionMode.Using = !state.SelectionMode.Using
	}

	if rl.IsKeyPressed(rl.KeyM) || rl.IsKeyPressed(rl.KeyEscape) {
		state.AppState.View = utils.PAUSED
	}

	zoomMult := float32(rl.GetMouseWheelMove()) * 0.1
	if state.Camera.Zoom+zoomMult > 0.2 {
		state.Camera.Zoom += zoomMult
	}

	if state.SelectionMode.Using {
		if state.Player.Turn.Actions > 0 {
			if rl.IsKeyPressed(rl.KeyB) {
				if tile, ok := getTile(state.SelectionMode.Pos); ok {
					if tile.Destroy() {
						state.Player.Turn.Actions--
					}
				}
			}
			if rl.IsKeyPressed(rl.KeyV) {
				for _, enemy := range state.Enemies {
					if enemy.Pos == state.SelectionMode.Pos {
						state.Player.Attack(enemy)
					}
				}
			}
		}
	}

	if utils.DebugMode {
		if rl.IsKeyPressed(rl.KeyF1) {
			state.DebugDisplay.Enabled = !state.DebugDisplay.Enabled
		}

		if rl.IsKeyPressed(rl.KeyI) {
			state.Player.Stats.Visibility++
		}

		if rl.IsKeyPressed(rl.KeyK) {
			state.Player.Stats.Visibility--
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		state.Player.Turn.Done = true
	}

	state.Camera.Target.X = float32(state.Player.Pos.X)
	state.Camera.Target.Y = float32(state.Player.Pos.Y)
}

func movePlayer(player *Player, tiles *[][]*Tile) {

	if player.Turn.Movement > 0 {
		p_x := player.Pos.X
		p_y := player.Pos.Y

		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			p_x -= TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			p_x += TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
			p_y -= TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
			p_y += TILE_SIZE
		}

		npos := utils.IVector2{X: p_x - PLAYER_OFFSET_X, Y: p_y - PLAYER_OFFSET_Y}

		if tile, ok := getTile(npos); ok {
			if tile.Block {
				return
			}
		} else {
			return
		}

		if p_x != player.Pos.X || p_y != player.Pos.Y {
			player.Pos.X = p_x
			player.Pos.Y = p_y
			player.Turn.Movement--
		}
	}
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
