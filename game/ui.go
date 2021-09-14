package game

import (
	"math"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIState struct {
	CharacterPanelOpen bool
	SelectionMode      SelectionMode
	DebugDisplay       DebugDisplayData
}

type SelectionMode struct {
	Using bool
	Pos   utils.IVector2
}

func NewUIState(player *Player) UIState {
	return UIState{
		CharacterPanelOpen: false,
		DebugDisplay: DebugDisplayData{
			Enabled:         false,
			TileDisplayMode: DD_TILE_NO_DISPLAY,
			TileLightFx:     true,
		},
		SelectionMode: SelectionMode{
			Using: false,
			Pos:   player.Pos,
		},
	}
}

func drawSelectionCursor() {
	if state.UIState.SelectionMode.Using {
		alpha := float32((math.Cos(3.0*float64(rl.GetTime())) + 1) * 0.5)
		rl.DrawTexture(*rendering.GetUISprite(rendering.SPRITE_SELECTION_MARK), state.UIState.SelectionMode.Pos.X, state.UIState.SelectionMode.Pos.Y, rl.ColorAlpha(rl.White, alpha))
	}
}

func drawUI() {
	RES := state.AppState.Settings.Resolution
	if !state.Player.Turn.Done {

		for h := 0; h < int(state.Player.Turn.Actions); h++ {
			rl.DrawTexture(*rendering.GetUISprite(rendering.SPRITE_ACTION_MARK), RES.X-100, int32(10+h*int(TILE_SIZE)+5), rl.White)
		}

		for m := 0; m < int(state.Player.Turn.Movement); m++ {
			rl.DrawTexture(*rendering.GetUISprite(rendering.SPRITE_MOVEMENT_MARK), RES.X-60, int32(10+m*int(TILE_SIZE)+5), rl.White)
		}

		if !(state.Player.Turn.Actions > 0) || !(state.Player.Turn.Movement > 0) {
			rendering.DrawMainText(rl.NewVector2(float32(RES.X/2), float32(RES.Y)/1.1), 48.0, "ENTER TO END TURN", rl.RayWhite)
		}
	} else {
		rendering.DrawMainText(rl.NewVector2(float32(RES.X/2), float32(RES.Y)/8.0), 48.0, "PROCESSING TURNS", rl.RayWhite)
	}

	if state.UIState.DebugDisplay.Enabled {
		drawDebugSettings()
		drawDebugInfo()
	}
}
