package rendering

import (
	"fmt"
	"utils"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Assets utils.RenderingAssets
var appState *utils.State

const (
	TILE_FLOOR_STONE    = iota
	TILE_WALL_STONE     = iota
	TILE_WALL_MOSS      = iota
	TILE_FLOOR_SPAWN    = iota
	TILE_FLOOR_OBS      = iota
	TILE_FLOOR_STONE_BL = iota
)

func LoadAssets(state *utils.State) *utils.RenderingAssets {
	missingImg := rl.GenImageColor(32, 32, rl.Pink)
	missingTexture := rl.LoadTextureFromImage(missingImg)
	rl.UnloadImage(missingImg)

	main, sec := loadFonts()
	Assets = utils.RenderingAssets{
		TileTextures:     loadTileTextures(),
		CharacterSprites: loadCharacterSprites(),
		UISprites:        loadUISprites(),
		MissingTexture:   &missingTexture,
		MainFont:         main,
		SecondaryFont:    sec,
		TestTextures:     buildTileSet("wall_stone_tile"),
	}
	loadGUIStylesheet()
	appState = state
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
	for _, t := range Assets.TestTextures {
		rl.UnloadTexture(t)
	}
	rl.UnloadFont(Assets.MainFont)
	rl.UnloadFont(Assets.SecondaryFont)
}

func loadTileTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 7)
	texturelist[TILE_FLOOR_STONE] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile.png"))
	texturelist[TILE_WALL_STONE] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_stone_tile.png"))
	texturelist[TILE_WALL_MOSS] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "wall_moss_tile.png"))
	texturelist[TILE_FLOOR_SPAWN] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_spawn_tile.png"))
	texturelist[TILE_FLOOR_OBS] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_obs_tile.png"))
	texturelist[TILE_FLOOR_STONE_BL] = rl.LoadTexture(utils.GetAssetPath(utils.TEXTURE, "floor_stone_tile_bl.png"))

	return texturelist
}

func buildTileSet(name string) []rl.Texture2D {
	base := rl.LoadImage(utils.GetAssetPath(utils.TEXTURE, fmt.Sprintf("%v.png", name)))
	rect := rl.NewRectangle(0.0, 0.0, 32.0, 32.0)

	up := rl.LoadImage(utils.GetAssetPath(utils.TEXTURE, fmt.Sprintf("%v_vert.png", name)))
	down := rl.ImageCopy(up)
	rl.ImageFlipVertical(down)

	right := rl.LoadImage(utils.GetAssetPath(utils.TEXTURE, fmt.Sprintf("%v_hor.png", name)))
	left := rl.ImageCopy(right)
	rl.ImageFlipHorizontal(left)

	up_right := rl.LoadImage(utils.GetAssetPath(utils.TEXTURE, fmt.Sprintf("%v_cor.png", name)))
	down_right := rl.ImageCopy(up_right)
	rl.ImageFlipVertical(down_right)
	up_left := rl.ImageCopy(up_right)
	rl.ImageFlipHorizontal(up_left)
	down_left := rl.ImageCopy(down_right)
	rl.ImageFlipHorizontal(down_left)

	up_tile := rl.ImageCopy(base)
	rl.ImageDraw(up_tile, up, rect, rect, rl.White)
	right_tile := rl.ImageCopy(base)
	rl.ImageDraw(right_tile, right, rect, rect, rl.White)
	down_tile := rl.ImageCopy(base)
	rl.ImageDraw(down_tile, down, rect, rect, rl.White)
	left_tile := rl.ImageCopy(base)
	rl.ImageDraw(left_tile, left, rect, rect, rl.White)
	up_right_tile := rl.ImageCopy(up_tile)
	rl.ImageDraw(up_right_tile, right, rect, rect, rl.White)
	rl.ImageDraw(up_right_tile, up_right, rect, rect, rl.White)
	up_left_tile := rl.ImageCopy(up_tile)
	rl.ImageDraw(up_left_tile, left, rect, rect, rl.White)
	rl.ImageDraw(up_left_tile, up_left, rect, rect, rl.White)
	up_down_tile := rl.ImageCopy(up_tile)
	rl.ImageDraw(up_down_tile, down, rect, rect, rl.White)
	down_right_tile := rl.ImageCopy(down_tile)
	rl.ImageDraw(down_right_tile, right, rect, rect, rl.White)
	rl.ImageDraw(down_right_tile, down_right, rect, rect, rl.White)
	down_left_tile := rl.ImageCopy(down_tile)
	rl.ImageDraw(down_left_tile, left, rect, rect, rl.White)
	rl.ImageDraw(down_left_tile, down_left, rect, rect, rl.White)
	right_left_tile := rl.ImageCopy(right_tile)
	rl.ImageDraw(right_left_tile, left, rect, rect, rl.White)
	up_right_down_tile := rl.ImageCopy(up_right_tile)
	rl.ImageDraw(up_right_down_tile, down, rect, rect, rl.White)
	rl.ImageDraw(up_right_down_tile, down_right, rect, rect, rl.White)
	up_right_left_tile := rl.ImageCopy(up_right_tile)
	rl.ImageDraw(up_right_left_tile, left, rect, rect, rl.White)
	rl.ImageDraw(up_right_left_tile, up_left, rect, rect, rl.White)
	up_down_left_tile := rl.ImageCopy(up_left_tile)
	rl.ImageDraw(up_down_left_tile, down, rect, rect, rl.White)
	rl.ImageDraw(up_down_left_tile, down_left, rect, rect, rl.White)
	right_down_left_tile := rl.ImageCopy(down_left_tile)
	rl.ImageDraw(right_down_left_tile, right, rect, rect, rl.White)
	rl.ImageDraw(right_down_left_tile, down_right, rect, rect, rl.White)
	all_tile := rl.ImageCopy(up_right_left_tile)
	rl.ImageDraw(all_tile, down, rect, rect, rl.White)
	rl.ImageDraw(all_tile, down_right, rect, rect, rl.White)
	rl.ImageDraw(all_tile, down_left, rect, rect, rl.White)

	texturelist := make([]rl.Texture2D, 16*2)
	texturelist[0] = rl.LoadTextureFromImage(base)
	texturelist[1] = rl.LoadTextureFromImage(up_tile)
	texturelist[2] = rl.LoadTextureFromImage(right_tile)
	texturelist[3] = rl.LoadTextureFromImage(up_right_tile)
	texturelist[4] = rl.LoadTextureFromImage(down_tile)
	texturelist[5] = rl.LoadTextureFromImage(up_down_tile)
	texturelist[6] = rl.LoadTextureFromImage(down_right_tile)
	texturelist[7] = rl.LoadTextureFromImage(up_right_down_tile)
	texturelist[8] = rl.LoadTextureFromImage(left_tile)
	texturelist[9] = rl.LoadTextureFromImage(up_left_tile)
	texturelist[10] = rl.LoadTextureFromImage(right_left_tile)
	texturelist[11] = rl.LoadTextureFromImage(up_right_left_tile)
	texturelist[12] = rl.LoadTextureFromImage(down_left_tile)
	texturelist[13] = rl.LoadTextureFromImage(up_down_left_tile)
	texturelist[14] = rl.LoadTextureFromImage(right_down_left_tile)
	texturelist[15] = rl.LoadTextureFromImage(all_tile)

	in_up_right := rl.LoadImage(utils.GetAssetPath(utils.TEXTURE, fmt.Sprintf("%v_incor.png", name)))
	in_down_right := rl.ImageCopy(in_up_right)
	rl.ImageFlipVertical(in_down_right)
	in_up_left := rl.ImageCopy(in_up_right)
	rl.ImageFlipHorizontal(in_up_left)
	in_down_left := rl.ImageCopy(in_down_right)
	rl.ImageFlipHorizontal(in_down_left)

	in_up_right_tile := rl.ImageCopy(base)
	rl.ImageDraw(in_up_right_tile, in_up_right, rect, rect, rl.White)
	in_down_right_tile := rl.ImageCopy(base)
	rl.ImageDraw(in_down_right_tile, in_down_right, rect, rect, rl.White)
	in_down_left_tile := rl.ImageCopy(base)
	rl.ImageDraw(in_down_left_tile, in_down_left, rect, rect, rl.White)
	in_up_left_tile := rl.ImageCopy(base)
	rl.ImageDraw(in_up_left_tile, in_up_left, rect, rect, rl.White)

	texturelist[16] = texturelist[0]
	texturelist[16+1] = rl.LoadTextureFromImage(in_up_right_tile)
	texturelist[16+2] = rl.LoadTextureFromImage(in_down_right_tile)
	texturelist[16+3] = rl.LoadTextureFromImage(base)
	texturelist[16+4] = rl.LoadTextureFromImage(in_down_left_tile)
	texturelist[16+5] = rl.LoadTextureFromImage(base)
	texturelist[16+6] = rl.LoadTextureFromImage(base)
	texturelist[16+7] = rl.LoadTextureFromImage(base)
	texturelist[16+8] = rl.LoadTextureFromImage(in_up_left_tile)
	texturelist[16+9] = rl.LoadTextureFromImage(base)
	texturelist[16+10] = rl.LoadTextureFromImage(base)
	texturelist[16+11] = rl.LoadTextureFromImage(base)
	texturelist[16+12] = rl.LoadTextureFromImage(base)
	texturelist[16+13] = rl.LoadTextureFromImage(base)
	texturelist[16+14] = rl.LoadTextureFromImage(base)
	texturelist[16+15] = rl.LoadTextureFromImage(base)

	rl.UnloadImage(base)
	rl.UnloadImage(up)
	rl.UnloadImage(down)
	rl.UnloadImage(right)
	rl.UnloadImage(left)
	rl.UnloadImage(up_tile)
	rl.UnloadImage(right_tile)
	rl.UnloadImage(down_tile)
	rl.UnloadImage(left_tile)
	rl.UnloadImage(up_right_tile)
	rl.UnloadImage(up_left_tile)
	rl.UnloadImage(up_down_tile)
	rl.UnloadImage(down_right_tile)
	rl.UnloadImage(down_left_tile)
	rl.UnloadImage(right_left_tile)
	rl.UnloadImage(up_right_down_tile)
	rl.UnloadImage(up_right_left_tile)
	rl.UnloadImage(up_down_left_tile)
	rl.UnloadImage(right_down_left_tile)
	rl.UnloadImage(all_tile)
	rl.UnloadImage(in_up_right)

	rl.UnloadImage(in_up_right_tile)

	return texturelist
}

func GetTile(tileType int) *rl.Texture2D {
	if tex := &Assets.TileTextures[tileType]; tex.Height != 0 {
		return tex
	} else {
		return Assets.MissingTexture
	}
}

const (
	PLAYER_IDLE = iota
	GOBLIN_IDLE = iota
)

func loadCharacterSprites() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 2)
	texturelist[PLAYER_IDLE] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "player_idle.png"))
	texturelist[GOBLIN_IDLE] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "goblin_idle.png"))

	return texturelist
}

func GetCharacterSprite(state int) *rl.Texture2D {
	if tex := &Assets.CharacterSprites[state]; tex.Height != 0 {
		return tex
	} else {
		return Assets.MissingTexture
	}
}

const (
	SPRITE_ACTION_MARK    = iota
	SPRITE_MOVEMENT_MARK  = iota
	SPRITE_SELECTION_MARK = iota
	SPRITE_CHECKMARK      = iota
	SPRITE_CROSS          = iota
)

func loadUISprites() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 5)
	texturelist[SPRITE_ACTION_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "action_mark.png"))
	texturelist[SPRITE_MOVEMENT_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "movement_mark.png"))
	texturelist[SPRITE_SELECTION_MARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "selection_mark.png"))
	texturelist[SPRITE_CHECKMARK] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "checkmark.png"))
	texturelist[SPRITE_CROSS] = rl.LoadTexture(utils.GetAssetPath(utils.SPRITE, "cross.png"))

	return texturelist
}

func GetUISprite(uiAsset int) *rl.Texture2D {
	if tex := &Assets.UISprites[uiAsset]; tex.Height != 0 {
		return tex
	} else {
		return Assets.MissingTexture
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
