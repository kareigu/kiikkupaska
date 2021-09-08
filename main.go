package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"game"
	"rendering"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var RES = game.IVector2{
	X: 800,
	Y: 600,
}

// Handle cli arguments
func init() {
	if len(os.Args) > 1 {
		args := os.Args[1:]

		width, err := strconv.Atoi(args[0])
		if err == nil {
			RES.X = int32(width)
		}

		height, err := strconv.Atoi(args[1])
		if err == nil {
			RES.Y = int32(height)
		}

		log.Printf("Resolution changed to: %dx%d\n", width, height)
	}
}

func main() {
	rl.InitWindow(RES.X, RES.Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	tile_textures := rendering.LoadTileTextures()
	character_textures := rendering.LoadCharacterTextures()

	game.InitGame(&RES, &character_textures, &tile_textures)

	for !rl.WindowShouldClose() {
		rl.SetWindowTitle(fmt.Sprintf("kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))
		game.GameUpdate()
	}

	for _, t := range tile_textures {
		rl.UnloadTexture(t)
	}
	for _, t := range character_textures {
		rl.UnloadTexture(t)
	}

	rl.CloseWindow()
}
