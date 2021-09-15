package rendering

import (
	"fmt"
	"log"
	"time"
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
	for _, t := range Assets.TestTextures.Textures {
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

func BuildTileSet(name string) utils.TileSet {
	t := time.Now()
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

	parts := [13]*rl.Image{
		base,
		up,
		up_right,
		right,
		down_right,
		down,
		down_left,
		left,
		up_left,
		in_up_right,
		in_down_right,
		in_down_left,
		in_up_left,
	}

	var texturelist [4096]rl.Texture2D
	for i := 0; i < 4096; i++ {
		texturelist[i] = rl.LoadTextureFromImage(buildTile(&parts, uint16(i)))
	}

	log.Printf("%v tileset built in %v", name, time.Since(t))

	return utils.TileSet{
		Parts:    parts,
		Textures: texturelist,
		Loaded:   true,
	}
}

func buildTile(parts *[13]*rl.Image, n uint16) *rl.Image {
	const TILE_SIZE float32 = 32.0
	tile := rl.ImageCopy(parts[0])
	rect := rl.NewRectangle(0.0, 0.0, TILE_SIZE, TILE_SIZE)
	//top_right := rl.NewRectangle(TILE_SIZE/2, 0.0, TILE_SIZE, TILE_SIZE/2)

	fmt.Printf("%v : %v\n", n, n&1)

	if n&1 > 0 {
		rl.ImageDraw(tile, parts[1], rect, rect, rl.White)
	}

	if n&8 > 0 {
		rl.ImageDraw(tile, parts[3], rect, rect, rl.White)
	}

	if n&64 > 0 {
		rl.ImageDraw(tile, parts[5], rect, rect, rl.White)
	}

	if n&512 > 0 {
		rl.ImageDraw(tile, parts[7], rect, rect, rl.White)
	}

	if n&2 > 0 {
		rl.ImageDraw(tile, parts[2], rect, rect, rl.White)
	}

	if n&4 > 0 {
		rl.ImageDraw(tile, parts[9], rect, rect, rl.White)
	}

	if n&16 > 0 {
		rl.ImageDraw(tile, parts[4], rect, rect, rl.White)
	}

	if n&32 > 0 {
		rl.ImageDraw(tile, parts[10], rect, rect, rl.White)
	}

	if n&128 > 0 {
		rl.ImageDraw(tile, parts[6], rect, rect, rl.White)
	}

	if n&256 > 0 {
		rl.ImageDraw(tile, parts[11], rect, rect, rl.White)
	}

	if n&1024 > 0 {
		rl.ImageDraw(tile, parts[8], rect, rect, rl.White)
	}

	if n&2048 > 0 {
		rl.ImageDraw(tile, parts[12], rect, rect, rl.White)
	}

	return tile
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
