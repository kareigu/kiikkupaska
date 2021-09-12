package utils

import (
	"log"
	"strconv"
	"strings"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAIN_MENU = iota
	PAUSED    = iota
	IN_GAME   = iota
)

const (
	TEXTURE    = iota
	SPRITE     = iota
	SOUND      = iota
	MUSIC      = iota
	STYLESHEET = iota
	FONT       = iota
)

const assetsFolder = "assets/"
const textureFolder = assetsFolder + "textures/"
const spriteFolder = assetsFolder + "sprites/"
const soundFolder = assetsFolder + "sound/"
const stylesFolder = assetsFolder + "stylesheets/"
const fontsFolder = assetsFolder + "fonts/"
const musicFolder = assetsFolder + "music/"

type State struct {
	Loading      bool
	View         int
	Settings     Settings
	RenderAssets *RenderingAssets
}

type Settings struct {
	PanelVisible       bool
	Music              bool
	Resolution         IVector2
	SelectedResolution int
}

type RenderingAssets struct {
	TileTextures     []rl.Texture2D
	CharacterSprites []rl.Texture2D
	UISprites        []rl.Texture2D
	MissingTexture   *rl.Texture2D
	MainFont         rl.Font
	SecondaryFont    rl.Font
}

type IVector2 struct {
	X int32
	Y int32
}

func NewIVector2(x int32, y int32) IVector2 {
	return IVector2{
		X: x,
		Y: y,
	}
}

func (ivec IVector2) ToVec2() rl.Vector2 {
	return rl.Vector2{X: float32(ivec.X), Y: float32(ivec.Y)}
}

var appState *State
var DebugMode bool

func InitUtils(state *State, debug bool) {
	appState = state
	DebugMode = debug
}

func DebugPrint(v interface{}) {
	if DebugMode {
		log.Print(v)
	}
}

func GetAssetPath(asset_type int, path string) string {
	switch asset_type {
	case TEXTURE:
		return textureFolder + path
	case SPRITE:
		return spriteFolder + path
	case SOUND:
		return soundFolder + path
	case MUSIC:
		return musicFolder + path
	case STYLESHEET:
		return stylesFolder + path
	case FONT:
		return fontsFolder + path
	default:
		return "Invalid asset type"
	}
}

var resolutionList = []string{
	"800x600",
	"1024x768",
	"1280x720",
	"1600x900",
	"1200x1000",
	"1920x1080",
	"2560x1440",
	"3840x2160",
}

func stringToRes(s string) IVector2 {
	split := strings.Split(s, "x")

	x_string := strings.Replace(split[0], "x", "", -1)
	var x int32 = 800
	if r, err := strconv.Atoi(x_string); err == nil {
		x = int32(r)
	}

	var y int32 = 600
	if r, err := strconv.Atoi(split[1]); err == nil {
		y = int32(r)
	}

	return NewIVector2(x, y)
}

func handleResolutionChange(newRes IVector2) bool {
	if newRes != appState.Settings.Resolution {
		appState.Settings.Resolution = newRes
		rl.SetWindowSize(int(appState.Settings.Resolution.X), int(appState.Settings.Resolution.Y))
		return true
	}
	return false
}

func DrawSettingsPanel() {
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
		appState.Settings.Resolution.ToVec2().X/2.0-245.0,
		appState.Settings.Resolution.ToVec2().Y/2.0-220.0,
		60.0,
		25.0,
	)

	appState.Settings.SelectedResolution = rgui.ToggleGroup(resolutionBackground, resolutionList, appState.Settings.SelectedResolution)
	if handleResolutionChange(stringToRes(resolutionList[appState.Settings.SelectedResolution])) {
		log.Print("Switched resolution to ", resolutionList[appState.Settings.SelectedResolution])
	}

	DrawSecondaryText(
		rl.NewVector2(appState.Settings.Resolution.ToVec2().X/2.0, appState.Settings.Resolution.ToVec2().Y/2.0+120.0),
		24.0,
		"Music",
		rl.RayWhite,
	)

	musicCheckboxBackground := rl.NewRectangle(
		appState.Settings.Resolution.ToVec2().X/2.0-25.0,
		appState.Settings.Resolution.ToVec2().Y/2.0+120.0,
		25.0,
		25.0,
	)
	appState.Settings.Music = rgui.CheckBox(musicCheckboxBackground, appState.Settings.Music)

	closeButtonBackground := rl.NewRectangle(
		appState.Settings.Resolution.ToVec2().X/2.0-40.0,
		appState.Settings.Resolution.ToVec2().Y/2.0+220.0,
		80.0,
		25.0,
	)
	if rgui.Button(closeButtonBackground, "Close") {
		appState.Settings.PanelVisible = false
	}
}

func DrawButton(pos rl.Vector2, text string) bool {
	return rgui.Button(rl.Rectangle{X: pos.X, Y: pos.Y, Width: 100.0, Height: 25.0}, text)
}

func DrawMainText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	rl.DrawTextEx(appState.RenderAssets.MainFont, text, pos, size, 1.0, colour)
}

func DrawSecondaryText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	rl.DrawTextEx(appState.RenderAssets.SecondaryFont, text, pos, size, 1.0, colour)
}
