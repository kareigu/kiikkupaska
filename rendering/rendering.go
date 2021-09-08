package rendering

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTileTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 7)
	texturelist[0] = rl.LoadTexture("missing_tile.png")
	texturelist[1] = rl.LoadTexture("floor_stone_tile.png")
	texturelist[2] = rl.LoadTexture("wall_stone_tile.png")
	texturelist[3] = rl.LoadTexture("wall_moss_tile.png")
	texturelist[4] = rl.LoadTexture("floor_spawn_tile.png")
	texturelist[5] = rl.LoadTexture("floor_obs_tile.png")
	texturelist[6] = rl.LoadTexture("floor_stone_tile_bl.png")

	return texturelist
}

func LoadCharacterTextures() []rl.Texture2D {
	texturelist := make([]rl.Texture2D, 1)
	texturelist[0] = rl.LoadTexture("player_idle.png")

	return texturelist
}
