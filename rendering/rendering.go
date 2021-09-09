package rendering

import (
	"utils"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTileTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 7)
	texturelist[0] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "missing_tile.png"))
	texturelist[1] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile.png"))
	texturelist[2] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_stone_tile.png"))
	texturelist[3] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_moss_tile.png"))
	texturelist[4] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_spawn_tile.png"))
	texturelist[5] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_obs_tile.png"))
	texturelist[6] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile_bl.png"))

	return texturelist
}

func LoadCharacterTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 1)
	texturelist[0] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "player_idle.png"))

	return texturelist
}

func LoadUISprites() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 3)
	texturelist[0] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "action_mark.png"))
	texturelist[1] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "movement_mark.png"))
	texturelist[2] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "selection_mark.png"))

	return texturelist
}

func LoadGUIStylesheet() {
	rgui.LoadGuiStyle(utils.GetAssetPath(utils.STYLESHEET, "zahnrad.style"))
}

func LoadFonts() (rl.Font, rl.Font) {
	mainFont := rl.LoadFont(utils.GetAssetPath(utils.FONT, "setback.png"))
	rl.GenTextureMipmaps(&mainFont.Texture)
	rl.SetTextureFilter(mainFont.Texture, rl.FilterPoint)
	secFont := rl.LoadFont(utils.GetAssetPath(utils.FONT, "alpha_beta.png"))
	rl.GenTextureMipmaps(&secFont.Texture)
	rl.SetTextureFilter(secFont.Texture, rl.FilterPoint)

	return mainFont, secFont
}
