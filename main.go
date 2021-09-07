package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	simplex "github.com/ojrac/opensimplex-go"
)

const TILE_SIZE int32 = 32
const PLAYER_OFFSET_X int32 = 5
const PLAYER_OFFSET_Y int32 = -2

var RES_X int32 = 800
var RES_Y int32 = 600

type IVector2 struct {
	X int32
	Y int32
}

type Tile struct {
	Colour rl.Color
	Pos    IVector2
	Block  bool
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

func charToTile(c string, pos IVector2) Tile {
	switch c {
	case "-":
		return Tile{
			Colour: rl.Black,
			Pos:    pos,
			Block:  false,
		}
	case "@":
		return Tile{
			Colour: rl.DarkGray,
			Pos:    pos,
			Block:  true,
		}
	case "!":
		return Tile{
			Colour: rl.DarkGreen,
			Pos:    pos,
			Block:  true,
		}
	case "b":
		return Tile{
			Colour: rl.Red,
			Pos:    pos,
			Block:  false,
		}
	case "_":
		return Tile{
			Colour: rl.Gray,
			Pos:    pos,
			Block:  false,
		}
	default:
		return Tile{
			Colour: rl.Blue,
			Pos:    pos,
			Block:  false,
		}
	}
}

type Player struct {
	Pos    IVector2
	Sprite string
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

var noisemapstring string = ""

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

func init() {
	source := rand.NewSource(time.Now().UnixMilli())
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
}

func main() {
	rl.InitWindow(RES_X, RES_Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

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
		Pos:    IVector2{X: PLAYER_OFFSET_X, Y: PLAYER_OFFSET_Y},
		Sprite: "@",
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
			tile := charToTile(char, pos)
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
		rl.ClearBackground(rl.DarkGray)

		for _, tile := range tiles {
			rl.DrawRectangle(tile.Pos.X, tile.Pos.Y, TILE_SIZE, TILE_SIZE, tile.Colour)
		}

		rl.DrawText(player.Sprite, player.Pos.X, player.Pos.Y, TILE_SIZE, rl.Red)

		rl.EndMode2D()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
