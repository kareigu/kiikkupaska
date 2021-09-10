package rendering

import (
	"utils"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Assets utils.RenderingAssets

const (
	TILE_MISSING        = iota
	TILE_FLOOR_STONE    = iota
	TILE_WALL_STONE     = iota
	TILE_WALL_MOSS      = iota
	TILE_FLOOR_SPAWN    = iota
	TILE_FLOOR_OBS      = iota
	TILE_FLOOR_STONE_BL = iota
)

func LoadAssets() *utils.RenderingAssets {
	main, sec := loadFonts()
	Assets = utils.RenderingAssets{
		TileTextures:     loadTileTextures(),
		CharacterSprites: loadCharacterSprites(),
		UISprites:        loadUISprites(),
		MainFont:         main,
		SecondaryFont:    sec,
	}
	loadGUIStylesheet()
	return &Assets
}

func Cleanup() {
	for _, t := range Assets.TileTextures {
		rl.UnloadTexture(t)
	}
	for _, t := range Assets.CharacterSprites {
		rl.UnloadTexture(t)
	}
	for _, t := range Assets.UISprites {
		rl.UnloadTexture(t)
	}
	rl.UnloadFont(Assets.MainFont)
	rl.UnloadFont(Assets.SecondaryFont)
}

func loadTileTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 7)
	texturelist[TILE_MISSING] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "missing_tile.png"))
	texturelist[TILE_FLOOR_STONE] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile.png"))
	texturelist[TILE_WALL_STONE] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_stone_tile.png"))
	texturelist[TILE_WALL_MOSS] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_moss_tile.png"))
	texturelist[TILE_FLOOR_SPAWN] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_spawn_tile.png"))
	texturelist[TILE_FLOOR_OBS] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_obs_tile.png"))
	texturelist[TILE_FLOOR_STONE_BL] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile_bl.png"))

	return texturelist
}

func GetTile(tileType int) *rl.Texture2D {
	if tex := &Assets.TileTextures[tileType]; tex != nil {
		return tex
	} else {
		return &Assets.TileTextures[TILE_MISSING]
	}
}

const (
	PLAYER_IDLE = iota
)

func loadCharacterSprites() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 1)
	texturelist[PLAYER_IDLE] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "player_idle.png"))

	return texturelist
}

func GetPlayerSprite(state int) *rl.Texture2D {
	if tex := &Assets.CharacterSprites[state]; tex != nil {
		return tex
	} else {
		return &Assets.TileTextures[TILE_MISSING]
	}
}

const (
	SPRITE_ACTION_MARK    = iota
	SPRITE_MOVEMENT_MARK  = iota
	SPRITE_SELECTION_MARK = iota
)

func loadUISprites() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 3)
	texturelist[SPRITE_ACTION_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "action_mark.png"))
	texturelist[SPRITE_MOVEMENT_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "movement_mark.png"))
	texturelist[SPRITE_SELECTION_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "selection_mark.png"))

	return texturelist
}

func GetUISprite(uiAsset int) *rl.Texture2D {
	if tex := &Assets.UISprites[uiAsset]; tex != nil {
		return tex
	} else {
		return &Assets.TileTextures[TILE_MISSING]
	}
}

func loadGUIStylesheet() {
	rgui.LoadGuiStyle(utils.GetAssetPath(utils.STYLESHEET, "zahnrad.style"))
}

func loadFonts() (rl.Font, rl.Font) {
	mainFont := rl.LoadFont(utils.GetAssetPath(utils.FONT, "setback.png"))
	rl.GenTextureMipmaps(&mainFont.Texture)
	rl.SetTextureFilter(mainFont.Texture, rl.FilterPoint)
	secFont := rl.LoadFont(utils.GetAssetPath(utils.FONT, "alpha_beta.png"))
	rl.GenTextureMipmaps(&secFont.Texture)
	rl.SetTextureFilter(secFont.Texture, rl.FilterPoint)

	return mainFont, secFont
}
