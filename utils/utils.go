package utils

import (
	"log"

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
	Loading       bool
	View          int
	RES           IVector2
	Music         bool
	MainFont      rl.Font
	SecondaryFont rl.Font
}

type IVector2 struct {
	X int32
	Y int32
}

/* type I_IVector2 interface {
	ToVec2() rl.Vector2
} */

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

func DrawButton(pos rl.Vector2, text string) bool {
	return rgui.Button(rl.Rectangle{X: pos.X, Y: pos.Y, Width: 100.0, Height: 25.0}, text)
}

func DrawMainText(pos rl.Vector2, size float32, text string, colour rl.Color) {
	rl.DrawTextEx(appState.MainFont, text, pos, size, 1.0, colour)
}
