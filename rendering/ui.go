package rendering

import (
	"log"
	"utils"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawMenuButtons(menu int, exitWindow *bool) {
	topButtonPos := rl.NewVector2(float32(appState.Settings.Resolution.X)/2.0, float32(appState.Settings.Resolution.Y)/2.0+50.0)
	botButtonPos := rl.NewVector2(float32(appState.Settings.Resolution.X)/2.0, float32(appState.Settings.Resolution.Y)/2.0+150.0)
	DrawMainText(rl.Vector2{X: float32(appState.Settings.Resolution.X / 2), Y: float32(appState.Settings.Resolution.Y / 6)}, 96.0, "KIIKKUPASKAA", rl.RayWhite)

	settings := DrawButton(rl.NewVector2(float32(appState.Settings.Resolution.X)/2.0, float32(appState.Settings.Resolution.Y)/2.0+100.0), "SETTINGS")

	if menu == utils.MAIN_MENU {
		start := DrawButton(topButtonPos, "START")
		exit := DrawButton(botButtonPos, "QUIT")

		if start {
			appState.View = utils.IN_GAME
		}

		if exit {
			*exitWindow = true
		}
	}

	if menu == utils.PAUSED {
		resume := DrawButton(topButtonPos, "RESUME")
		exit := DrawButton(botButtonPos, "EXIT TO MENU")

		if resume {
			appState.View = utils.IN_GAME
		}

		if exit {
			appState.View = utils.MAIN_MENU
		}
	}

	if appState.Settings.PanelVisible {
		DrawSettingsPanel()
	}
	if settings {
		appState.Settings.PanelVisible = !appState.Settings.PanelVisible
	}
}

func DrawSettingsPanel() {
	appState.Settings.SelectedResolution = 0
	for i, res := range utils.ResolutionList {
		if res == utils.ResToString(appState.Settings.Resolution) {
			appState.Settings.SelectedResolution = i
		}
	}

	background := rl.NewRectangle(
		appState.Settings.Resolution.ToVec2().X/2.0-300.0,
		appState.Settings.Resolution.ToVec2().Y/2.0-250.0,
		600.0,
		500.0,
	)
	rl.DrawRectangleRounded(background, 0.05, 2, rl.NewColor(19, 26, 40, 255))

	DrawSecondaryText(
		rl.NewVector2(appState.Settings.Resolution.ToVec2().X/2.0, appState.Settings.Resolution.ToVec2().Y/2.0-250.0),
		24.0,
		"Resolution",
		rl.RayWhite,
	)

	resolutionBackground := rl.NewRectangle(
		appState.Settings.Resolution.ToVec2().X/2.0-280.0,
		appState.Settings.Resolution.ToVec2().Y/2.0-220.0,
		60.0,
		25.0,
	)

	appState.Settings.SelectedResolution = rgui.ToggleGroup(resolutionBackground, utils.ResolutionList, appState.Settings.SelectedResolution)
	if utils.HandleResolutionChange(utils.StringToRes(utils.ResolutionList[appState.Settings.SelectedResolution])) {
		log.Print("Switched resolution to ", utils.ResolutionList[appState.Settings.SelectedResolution])
	}

	DrawSecondaryText(
		rl.NewVector2(
			appState.Settings.Resolution.ToVec2().X/2.0-25.0,
			appState.Settings.Resolution.ToVec2().Y/2.0+120.0,
		),
		25.0,
		"Music",
		rl.RayWhite,
	)

	musicCheckboxBackground := rl.NewRectangle(
		appState.Settings.Resolution.ToVec2().X/2.0+25.0,
		appState.Settings.Resolution.ToVec2().Y/2.0+120.0,
		25.0,
		25.0,
	)
	musicToggle := rgui.CheckBox(musicCheckboxBackground, appState.Settings.Music)
	if musicToggle != appState.Settings.Music {
		appState.Settings.Music = musicToggle
		utils.SaveSettingsFile()
	}

	checkboxTex := SPRITE_CROSS
	if appState.Settings.Music {
		checkboxTex = SPRITE_CHECKMARK
	}

	rl.DrawTextureV(
		appState.RenderAssets.UISprites[checkboxTex],
		rl.NewVector2(musicCheckboxBackground.X, musicCheckboxBackground.Y),
		rl.White,
	)

	closeButtonPos := rl.NewVector2(
		appState.Settings.Resolution.ToVec2().X/2.0,
		appState.Settings.Resolution.ToVec2().Y/2.0+220.0,
	)
	if DrawButton(closeButtonPos, "Close") {
		appState.Settings.PanelVisible = false
	}
}

func DrawButton(pos rl.Vector2, text string) bool {
	const width = 100.0
	const height = 25.0
	const textPadding = 4

	pos.X -= width / 2.0
	textHeight := appState.RenderAssets.SecondaryFont.BaseSize
	textWidth := rl.MeasureText(text, textHeight)
	bounds := rl.NewRectangle(pos.X, pos.Y, width, height)

	rgui.ConstrainRectangle(&bounds, textWidth, textWidth+textPadding, textHeight, textHeight+textPadding/2)

	state := rgui.GetInteractionState(bounds)
	base_colour := rl.NewColor(44, 60, 92, 255)
	base_border_colour := rl.NewColor(193, 153, 33, 255)
	focus_colour := rl.NewColor(61, 83, 128, 255)
	focus_border_colour := rl.NewColor(124, 118, 101, 255)

	colour := base_colour
	border_colour := base_border_colour

	if state == rgui.Focused {
		colour = focus_colour
		border_colour = focus_border_colour
	}
	if state == rgui.Clicked {
		colour = base_border_colour
		border_colour = base_colour
	}

	b := bounds.ToInt32()
	rgui.DrawBorderedRectangle(b, 2, border_colour, colour)
	textPos := rl.NewVector2(
		float32(b.X+(b.Width/2)+textPadding),
		float32(b.Y+((b.Height/2)-(textHeight/2))),
	)
	DrawSecondaryText(textPos, float32(textHeight), text, rl.RayWhite)

	return state == rgui.Clicked
}

func DrawDefaultText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	width := rl.MeasureText(text, int32(size))
	pos.X -= float32(width) / 2.0
	rl.DrawText(text, int32(pos.X), int32(pos.Y), int32(size), colour)
}

func DrawMainText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	width := rl.MeasureText(text, int32(size))
	pos.X -= float32(width / 2)
	rl.DrawTextEx(appState.RenderAssets.MainFont, text, pos, size, 1.0, colour)
}

func DrawSecondaryText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	width := rl.MeasureText(text, int32(size))
	pos.X -= float32(width) / 2.0
	rl.DrawTextEx(appState.RenderAssets.SecondaryFont, text, pos, size, 1.0, colour)
}
