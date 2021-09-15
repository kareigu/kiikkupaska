package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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
const tilesFolder = assetsFolder + "tiles/"
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
	TestTextures     TileSet
}

type TileSet struct {
	Parts    [13]*rl.Image
	Textures [4096]rl.Texture2D
	Loaded   bool
}

func (tileset *TileSet) GetTexture(neighbours uint16) *rl.Texture2D {
	if neighbours > 4096 {
		return appState.RenderAssets.MissingTexture
	}

	if tex := tileset.Textures[neighbours]; tex.Height != 0 {
		return &tex
	} else {
		return appState.RenderAssets.MissingTexture
	}
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
		return tilesFolder + path
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

var ResolutionList = []string{
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

func StringToRes(s string) IVector2 {
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

func ResToString(res IVector2) string {
	return fmt.Sprintf("%dx%d", res.X, res.Y)
}

func HandleResolutionChange(newRes IVector2) bool {
	if newRes != appState.Settings.Resolution {
		appState.Settings.Resolution = newRes
		rl.SetWindowSize(int(appState.Settings.Resolution.X), int(appState.Settings.Resolution.Y))
		centerWindow()
		SaveSettingsFile()
		return true
	}
	return false
}

func centerWindow() {
	currMonitor := rl.GetCurrentMonitor()
	monitorRes := NewIVector2(int32(rl.GetMonitorWidth(currMonitor)), int32(rl.GetMonitorHeight(currMonitor)))
	windowRes := appState.Settings.Resolution
	diff := rl.Vector2Subtract(monitorRes.ToVec2(), windowRes.ToVec2())
	diff = rl.Vector2DivideV(diff, rl.NewVector2(2.0, 2.0))
	rl.SetWindowPosition(int(diff.X), int(diff.Y))
}

func SaveSettingsFile() {
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
			SaveSettingsFile()
			loadSettingsFile(overrideRes)
		} else {
			appState.Settings.Music = settings.Music
			if !overrideRes {
				newRes := NewIVector2(int32(settings.ResolutionWidth), int32(settings.ResolutionHeight))
				appState.Settings.Resolution = newRes
			}
		}
	} else {
		log.Println("Settings file missing, writing with default settings")
		SaveSettingsFile()
		loadSettingsFile(overrideRes)
	}
}
