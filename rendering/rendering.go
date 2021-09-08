package rendering

import (
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
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

func LoadTileTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 7)
	texturelist[0] = rl.LoadTexture(getAssetPath(TEXTURE, "missing_tile.png"))
	texturelist[1] = rl.LoadTexture(getAssetPath(TEXTURE, "floor_stone_tile.png"))
	texturelist[2] = rl.LoadTexture(getAssetPath(TEXTURE, "wall_stone_tile.png"))
	texturelist[3] = rl.LoadTexture(getAssetPath(TEXTURE, "wall_moss_tile.png"))
	texturelist[4] = rl.LoadTexture(getAssetPath(TEXTURE, "floor_spawn_tile.png"))
	texturelist[5] = rl.LoadTexture(getAssetPath(TEXTURE, "floor_obs_tile.png"))
	texturelist[6] = rl.LoadTexture(getAssetPath(TEXTURE, "floor_stone_tile_bl.png"))

	return texturelist
}

func LoadCharacterTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 1)
	texturelist[0] = rl.LoadTexture(getAssetPath(SPRITE, "player_idle.png"))

	return texturelist
}

func LoadGUIStylesheet() {
	raygui.LoadGuiStyle(getAssetPath(STYLESHEET, "zahnrad.style"))
}

func LoadFont() rl.Font {
	font := rl.LoadFont(getAssetPath(FONT, "setback.png"))
	rl.GenTextureMipmaps(&font.Texture)
	rl.SetTextureFilter(font.Texture, rl.FilterPoint)
	return font
}

func getAssetPath(asset_type int, path string) string {
	switch asset_type {
	case TEXTURE:
		return textureFolder + path
	case SPRITE:
		return spriteFolder + path
	case SOUND:
		return soundFolder + path
	case MUSIC:
		return "not_implemented"
	case STYLESHEET:
		return stylesFolder + path
	case FONT:
		return fontsFolder + path
	default:
		return "Invalid asset type"
	}
}
