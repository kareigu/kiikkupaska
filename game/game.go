package game

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	simplex "github.com/ojrac/opensimplex-go"
)

const TILE_SIZE int32 = 32
const PLAYER_OFFSET_X int32 = 0
const PLAYER_OFFSET_Y int32 = 0

type GameState struct {
	Resolution *IVector2
	Camera     *rl.Camera2D
	Player     *Player
	Map        *map[IVector2]Tile
}

type IVector2 struct {
	X int32
	Y int32
}

type Tile struct {
	Texture rl.Texture2D
	Pos     IVector2
	Block   bool
}

type Player struct {
	Pos        IVector2
	Sprite     rl.Texture2D
	Visibility int8
}

var state GameState

func InitGame(res *IVector2, character_textures *[]rl.Texture2D, tile_textures *[]rl.Texture2D) *GameState {
	player, cam := initPlayerAndCam((*res), character_textures)
	state = GameState{
		Resolution: res,
		Player:     player,
		Camera:     cam,
		Map:        nil,
	}
	state.Map = generateTiles(tile_textures)
	return &state
}

func GameUpdate() {
	movePlayer(state.Player, state.Map, rl.GetKeyPressed())

	if rl.IsKeyPressed(rl.KeyI) {
		state.Player.Visibility++
	}

	if rl.IsKeyPressed(rl.KeyK) {
		state.Player.Visibility--
	}

	state.Camera.Target.X = float32(state.Player.Pos.X)
	state.Camera.Target.Y = float32(state.Player.Pos.Y)

	rl.BeginDrawing()
	rl.BeginMode2D(*state.Camera)
	rl.ClearBackground(rl.Black)

	for _, tile := range *state.Map {
		// Check if tile coordinates are in player visibility range
		// If not don't bother rendering it
		if colour, ok := checkTileVisibility(state.Player, &tile); ok {
			rl.DrawTexture(tile.Texture, tile.Pos.X, tile.Pos.Y, colour)

			//! Tile light debug display
			//rl.DrawText(fmt.Sprintf("%d", colour.A), tile.Pos.X, tile.Pos.Y, 12, rl.Red)
			//! Tile distance debug display
			/* dist := getTileDistanceToPlayer(&player, &tile)
			rl.DrawText(fmt.Sprintf("%.1f", math.Min(float64(dist), 10.0)), tile.Pos.X, tile.Pos.Y, 12, rl.Red) */
		}
	}

	rl.DrawTexture(state.Player.Sprite, state.Player.Pos.X, state.Player.Pos.Y, rl.White)
	//rl.DrawText(player.Sprite, player.Pos.X, player.Pos.Y, TILE_SIZE, rl.Red)

	rl.EndMode2D()
	rl.EndDrawing()
}

func initPlayerAndCam(res IVector2, character_textures *[]rl.Texture2D) (*Player, *rl.Camera2D) {
	cam := rl.Camera2D{
		Offset: rl.Vector2{
			X: float32(res.X / 2),
			Y: float32(res.Y / 2),
		},
		Target: rl.Vector2{
			X: 0.0,
			Y: 0.0,
		},
		Rotation: 0.0,
		Zoom:     1.0,
	}

	player := Player{
		Pos:        IVector2{X: PLAYER_OFFSET_X, Y: PLAYER_OFFSET_Y},
		Sprite:     (*character_textures)[0],
		Visibility: 8,
	}

	return &player, &cam
}

func generateTiles(tile_textures *[]rl.Texture2D) *map[IVector2]Tile {
	mapstring := generateMap()
	tiles := make(map[IVector2]Tile)
	player := state.Player

	for y, row := range strings.Split(mapstring, "\n") {
		for x, char := range strings.Split(row, "") {
			pos_x := int32(x) * TILE_SIZE
			pos_y := int32(y) * TILE_SIZE
			if char == "P" {
				if player.Pos.X == PLAYER_OFFSET_X && player.Pos.Y == PLAYER_OFFSET_Y {
					player.Pos.X = pos_x + PLAYER_OFFSET_X
					player.Pos.Y = pos_y + PLAYER_OFFSET_Y
				} else {
					if rand.Float32() < 0.1 {
						player.Pos.X = pos_x + PLAYER_OFFSET_X
						player.Pos.Y = pos_y + PLAYER_OFFSET_Y
					}
				}
			}
			pos := IVector2{X: pos_x, Y: pos_y}
			tile := charToTile(tile_textures, char, pos)
			tiles[pos] = tile
		}
	}
	log.Print(mapstring)
	return &tiles
}

func checkTileVisibility(player *Player, tile *Tile) (rl.Color, bool) {
	visrange := int32(player.Visibility) * TILE_SIZE
	if tile.Pos.X > player.Pos.X+(visrange) || tile.Pos.X < player.Pos.X-(visrange) || tile.Pos.Y > player.Pos.Y+(visrange) || tile.Pos.Y < player.Pos.Y-(visrange) {
		return rl.Black, false
	} else {
		distance_alpha := float32(getTileDistanceToPlayer(player, tile)) / float32(player.Visibility)
		colour := rl.ColorAlpha(rl.White, distance_alpha)
		// Reverse alpha to make closer tiles brighter instead of darker
		colour.A = uint8(math.Abs(float64(colour.A) - 255.0))
		return colour, true
	}
}

func getTileDistanceToPlayer(player *Player, tile *Tile) float32 {
	tile_vec := IVec2ToVec2(tile.Pos)
	player_vec := IVec2ToVec2(player.Pos)
	distance := rl.Vector2Distance(tile_vec, player_vec) / float32(TILE_SIZE)

	return distance
}

func IVec2ToVec2(ivec IVector2) rl.Vector2 {
	return rl.Vector2{X: float32(ivec.X), Y: float32(ivec.Y)}
}

func reverseRange(v float32) float32 {
	return float32(math.Abs((float64(v) - 1.0) / (0.0 - 1.0)))
}

func generateMap() string {
	mapstring := ""
	t := time.Now()
	log.Println("Map generation started")
	source := rand.NewSource(t.UnixMilli())
	rng := rand.New(source)
	noise := simplex.New(rng.Int63())
	var gen_i float64 = 0.0
	var gen_j float64 = 0.0

	for gen_i <= 10.0 {
		for gen_j <= 10.0 {
			if gen_i == 0.0 || gen_i > 9.9 || gen_j == 0.0 || gen_j > 9.9 {
				mapstring += "@"
			} else {
				val := noise.Eval2(gen_i, gen_j)
				if val > 0.7 || val < -0.7 {
					mapstring += "-"
				} else if val > 0.1 {
					mapstring += "@"
				} else if val > -0.6 {
					if rand.Float32() < 0.01 {
						mapstring += "P"
					} else {
						mapstring += "_"
					}
				} else {
					mapstring += "!"
				}
			}

			gen_j += 0.1
		}
		mapstring += "\n"
		gen_j = 0.0
		gen_i += 0.1
	}

	log.Println("Map generation finished in ", time.Since(t))
	return mapstring
}

func charToTile(texturelist *[]rl.Texture2D, c string, pos IVector2) Tile {
	switch c {
	case "@":
		return Tile{
			Texture: (*texturelist)[2],
			Pos:     pos,
			Block:   true,
		}
	case "_":
		ti := 1
		if rand.Float32() < 0.01 {
			ti = 6
		}
		return Tile{
			Texture: (*texturelist)[ti],
			Pos:     pos,
			Block:   false,
		}
	case "!":
		return Tile{
			Texture: (*texturelist)[3],
			Pos:     pos,
			Block:   true,
		}
	case "P":
		return Tile{
			Texture: (*texturelist)[4],
			Pos:     pos,
			Block:   false,
		}
	case "-":
		return Tile{
			Texture: (*texturelist)[5],
			Pos:     pos,
			Block:   false,
		}
	default:
		return Tile{
			Texture: (*texturelist)[0],
			Pos:     pos,
			Block:   false,
		}
	}
}

func movePlayer(player *Player, tiles *map[IVector2]Tile, key int32) {
	p_x := player.Pos.X
	p_y := player.Pos.Y

	if rl.IsKeyPressed(rl.KeyLeft) {
		p_x -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		p_x += TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		p_y -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		p_y += TILE_SIZE
	}

	npos := IVector2{X: p_x - PLAYER_OFFSET_X, Y: p_y - PLAYER_OFFSET_Y}

	if tile, ok := (*tiles)[npos]; ok {
		if tile.Block {
			return
		}
	} else {
		return
	}

	player.Pos.X = p_x
	player.Pos.Y = p_y
}
