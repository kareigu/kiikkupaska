package utils

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
	Loading bool
	View    int
	RES     IVector2
	Music   bool
}

type IVector2 struct {
	X int32
	Y int32
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
