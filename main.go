package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	simplex "github.com/ojrac/opensimplex-go"
)

const TILE_SIZE int32 = 32
const PLAYER_OFFSET_X int32 = 0
const PLAYER_OFFSET_Y int32 = 0

var RES_X int32 = 800
var RES_Y int32 = 600

type IVector2 struct {
	X int32
	Y int32
}

type Tile struct {
	Texture rl.Texture2D
	Pos     IVector2
	Block   bool
}

const mapstring string = `mapstring
--------------------------------
-----@@@@@@@@@@@@@@@@@@@@-------
-----@@@@@@@@@@@@@@@@@@@@-------
-----@@________!_______@@-------
-----@@________!_______@@-------
-----@@________!_______@@-------
-----@@_______!!!______@@-------
-----@@________________@@-------
-----@@____P___________@@-------
-----@@________________@@-------
-----@@________________@@-------
-----@@@@@@@@@@@@@@@@@@@@-------
-----@@@@@@@@@@@@@@@@@@@@-------
--------------------------------
`

func randColour() rl.Color {
	switch rand.Intn(5) {
	case 0:
		return rl.Black
	case 1:
		return rl.White
	case 2:
		return rl.Green
	case 3:
		return rl.Red
	default:
		return rl.Blue
	}
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

func LoadTileTextures(texturelist *[]rl.Texture2D) {
	*texturelist = make([]rl.Texture2D, 7)
	(*texturelist)[0] = rl.LoadTexture("missing_tile.png")
	(*texturelist)[1] = rl.LoadTexture("floor_stone_tile.png")
	(*texturelist)[2] = rl.LoadTexture("wall_stone_tile.png")
	(*texturelist)[3] = rl.LoadTexture("wall_moss_tile.png")
	(*texturelist)[4] = rl.LoadTexture("floor_spawn_tile.png")
	(*texturelist)[5] = rl.LoadTexture("floor_obs_tile.png")
	(*texturelist)[6] = rl.LoadTexture("floor_stone_tile_bl.png")
}

type Player struct {
	Pos        IVector2
	Sprite     rl.Texture2D
	Visibility int8
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

func init() {
	if len(os.Args) > 1 {
		args := os.Args[1:]

		width, err := strconv.Atoi(args[0])
		if err == nil {
			RES_X = int32(width)
		}

		height, err := strconv.Atoi(args[1])
		if err == nil {
			RES_Y = int32(height)
		}

		log.Printf("Resolution changed to: %dx%d\n", width, height)
	}
}

var noisemapstring string = ""

func init() {
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
				noisemapstring += "@"
			} else {
				val := noise.Eval2(gen_i, gen_j)
				if val > 0.7 || val < -0.7 {
					noisemapstring += "-"
				} else if val > 0.1 {
					noisemapstring += "@"
				} else if val > -0.6 {
					if rand.Float32() < 0.01 {
						noisemapstring += "P"
					} else {
						noisemapstring += "_"
					}
				} else {
					noisemapstring += "!"
				}
			}

			gen_j += 0.1
		}
		noisemapstring += "\n"
		gen_j = 0.0
		gen_i += 0.1
	}

	log.Println("Map generation finished in ", time.Since(t))
}

var tile_textures []rl.Texture2D

func main() {
	rl.InitWindow(RES_X, RES_Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	LoadTileTextures(&tile_textures)

	cam := rl.Camera2D{
		Offset: rl.Vector2{
			X: float32(RES_X / 2),
			Y: float32(RES_Y / 2),
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
		Sprite:     rl.LoadTexture("player_idle.png"),
		Visibility: 8,
	}

	log.Print(noisemapstring)

	var tiles map[IVector2]Tile
	tiles = make(map[IVector2]Tile)

	for y, row := range strings.Split(noisemapstring, "\n") {
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
			tile := charToTile(&tile_textures, char, pos)
			tiles[pos] = tile
		}
	}

	for !rl.WindowShouldClose() {
		rl.SetWindowTitle(fmt.Sprintf("kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))

		movePlayer(&player, &tiles, rl.GetKeyPressed())

		cam.Target.X = float32(player.Pos.X)
		cam.Target.Y = float32(player.Pos.Y)

		rl.BeginDrawing()
		rl.BeginMode2D(cam)
		rl.ClearBackground(rl.Black)

		for _, tile := range tiles {
			if colour, ok := checkTileVisibility(&player, &tile); ok {
				rl.DrawTexture(tile.Texture, tile.Pos.X, tile.Pos.Y, colour)
			}
		}

		rl.DrawTexture(player.Sprite, player.Pos.X, player.Pos.Y, rl.White)
		//rl.DrawText(player.Sprite, player.Pos.X, player.Pos.Y, TILE_SIZE, rl.Red)

		rl.EndMode2D()
		rl.EndDrawing()
	}

	for _, t := range tile_textures {
		rl.UnloadTexture(t)
	}
	rl.UnloadTexture(player.Sprite)

	rl.CloseWindow()
}

func checkTileVisibility(player *Player, tile *Tile) (rl.Color, bool) {
	visrange := int32(player.Visibility) * TILE_SIZE
	if tile.Pos.X > player.Pos.X+(visrange) || tile.Pos.X < player.Pos.X-(visrange) || tile.Pos.Y > player.Pos.Y+(visrange) || tile.Pos.Y < player.Pos.Y-(visrange) {
		return rl.Black, false
	} else {
		distance_alpha := getTileDistanceToPlayer(player, tile)/float32(player.Visibility) - 0.15
		return rl.ColorAlpha(rl.White, reverseRange(distance_alpha)), true
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
