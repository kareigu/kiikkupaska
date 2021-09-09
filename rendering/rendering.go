package rendering

import (
	"utils"

	"github.com/gen2brain/raylib-go/raygui"
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

func LoadGUIStylesheet() {
	raygui.LoadGuiStyle(utils.GetAssetPath(utils.STYLESHEET, "zahnrad.style"))
}

func LoadFont() rl.Font {
	font := rl.LoadFont(utils.GetAssetPath(utils.FONT, "setback.png"))
	rl.GenTextureMipmaps(&font.Texture)
	rl.SetTextureFilter(font.Texture, rl.FilterPoint)
	return font
}
