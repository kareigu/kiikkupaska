package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	defaultRes := IVector2{
		X: 800,
		Y: 600,
	}
	appState = state
	DebugMode = debug
	loadSettingsFile(state.Settings.Resolution != defaultRes)
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

type Settings struct {
	PanelVisible       bool
	Music              bool
	Resolution         IVector2
	SelectedResolution int
}

type SettingsFile struct {
	Music            bool `json:"music"`
	ResolutionWidth  int  `json:"resolutionWidth"`
	ResolutionHeight int  `json:"resolutionHeight"`
}

var resolutionList = []string{
	"Custom",
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
	if s == "Custom" {
		return appState.Settings.Resolution
	} else {
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
}

func resToString(res IVector2) string {
	return fmt.Sprintf("%dx%d", res.X, res.Y)
}

func handleResolutionChange(newRes IVector2) bool {
	if newRes != appState.Settings.Resolution {
		appState.Settings.Resolution = newRes
		rl.SetWindowSize(int(appState.Settings.Resolution.X), int(appState.Settings.Resolution.Y))
		saveSettingsFile()
		return true
	}
	return false
}

func saveSettingsFile() {
	settings := SettingsFile{
		Music:            appState.Settings.Music,
		ResolutionWidth:  int(appState.Settings.Resolution.X),
		ResolutionHeight: int(appState.Settings.Resolution.Y),
	}

	file, _ := json.MarshalIndent(settings, "", "	")

	if f, err := os.OpenFile("settings.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		log.Println("Couldn't write settings file")
	} else {
		f.Write(file)
		log.Println("Rewrote settings file")
	}
}

func loadSettingsFile(overrideRes bool) {
	var settings SettingsFile

	if file, err := ioutil.ReadFile("settings.json"); err == nil {
		if err = json.Unmarshal(file, &settings); err != nil {
			log.Println("Malformed settings file, rewriting with default settings")
			saveSettingsFile()
			loadSettingsFile(overrideRes)
		} else {
			appState.Settings.Music = settings.Music
			if !overrideRes {
				newRes := NewIVector2(int32(settings.ResolutionWidth), int32(settings.ResolutionHeight))
				appState.Settings.Resolution = newRes
				rl.SetWindowSize(int(appState.Settings.Resolution.X), int(appState.Settings.Resolution.Y))
			}
		}
	} else {
		log.Println("Settings file missing, writing with default settings")
		saveSettingsFile()
		loadSettingsFile(overrideRes)
	}
}

func DrawSettingsPanel() {
	appState.Settings.SelectedResolution = 0
	for i, res := range resolutionList {
		if res == resToString(appState.Settings.Resolution) {
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

	appState.Settings.SelectedResolution = rgui.ToggleGroup(resolutionBackground, resolutionList, appState.Settings.SelectedResolution)
	if handleResolutionChange(stringToRes(resolutionList[appState.Settings.SelectedResolution])) {
		log.Print("Switched resolution to ", resolutionList[appState.Settings.SelectedResolution])
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
		saveSettingsFile()
	}

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
	const width = 100.0
	pos.X -= width / 2.0
	return rgui.Button(rl.Rectangle{X: pos.X, Y: pos.Y, Width: width, Height: 25.0}, text)
}

func DrawMainText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	width := rl.MeasureText(text, int32(size))
	pos.X -= float32(width) / 2.0
	rl.DrawTextEx(appState.RenderAssets.MainFont, text, pos, size, 1.0, colour)
}

func DrawSecondaryText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	width := rl.MeasureText(text, int32(size))
	pos.X -= float32(width) / 2.0
	rl.DrawTextEx(appState.RenderAssets.SecondaryFont, text, pos, size, 1.0, colour)
}
