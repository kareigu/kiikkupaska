package game

import (
	"fmt"
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

func drawCharacterPanel() {
	xPos := float32(25.0)
	yPos := state.AppState.Settings.Resolution.ToVec2().Y/2.0 - 250.0
	panelWidth := float32(400.0)
	panelHeight := float32(300.0)
	background := rl.NewRectangle(
		xPos,
		yPos,
		panelWidth,
		panelHeight,
	)

	rl.DrawRectangleRounded(background, 0.05, 2, rendering.PanelBackground)
	rl.DrawRectangleRoundedLines(background, 0.05, 2, 2.0, rendering.GoldAccent)

	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/2.0,
			yPos+8.0,
		),
		24.0,
		"CHARACTER",
		rl.RayWhite,
	)

	rl.DrawTexture(*rendering.GetCharacterSprite(state.Player.State), int32(xPos)+10, int32(yPos)+8, rl.White)

	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/4.0,
			yPos+32.0,
		),
		24.0,
		"STATS",
		rl.RayWhite,
	)

	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/5.5,
			yPos+64.0,
		),
		24.0,
		fmt.Sprintf("strength %v", state.Player.Stats.Strength),
		rl.RayWhite,
	)
	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/5.5,
			yPos+86.0,
		),
		24.0,
		fmt.Sprintf("dexterity %v", state.Player.Stats.Dexterity),
		rl.RayWhite,
	)
	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/5.5,
			yPos+108.0,
		),
		24.0,
		fmt.Sprintf("vitality %v", state.Player.Stats.Vitality),
		rl.RayWhite,
	)
	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/5.5,
			yPos+130.0,
		),
		24.0,
		fmt.Sprintf("movement %v", state.Player.Stats.Movement),
		rl.RayWhite,
	)
	rendering.DrawSecondaryText(
		rl.NewVector2(
			xPos+panelWidth/5.5,
			yPos+152.0,
		),
		24.0,
		fmt.Sprintf("visibility %v", state.Player.Stats.Visibility),
		rl.RayWhite,
	)
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

	if state.UIState.CharacterPanelOpen {
		drawCharacterPanel()
	}

	if state.UIState.DebugDisplay.Enabled {
		drawDebugSettings()
		drawDebugInfo()
	}
}
