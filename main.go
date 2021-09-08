package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"game"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var state utils.State

// Handle cli arguments
func init() {
	var x int32
	var y int32
	if len(os.Args) > 1 {
		args := os.Args[1:]

		width, err := strconv.Atoi(args[0])
		if err == nil {
			x = int32(width)
		}

		height, err := strconv.Atoi(args[1])
		if err == nil {
			y = int32(height)
		}

		log.Printf("Resolution changed to: %dx%d\n", width, height)
	}

	state = utils.State{
		Loading: false,
		View:    utils.MAIN_MENU,
		RES:     utils.IVector2{X: x, Y: y},
	}
}

func main() {
	rl.InitWindow(state.RES.X, state.RES.Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	tile_textures := rendering.LoadTileTextures()
	character_textures := rendering.LoadCharacterTextures()

	var gameState *game.GameState

	for !rl.WindowShouldClose() {
		rl.SetWindowTitle(fmt.Sprintf("kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))
		switch state.View {
		case utils.MAIN_MENU:
			if rl.IsKeyPressed(rl.KeyEnter) {
				state.View = utils.IN_GAME
			}

			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)
			rl.DrawText("Main Menu", state.RES.X/2, state.RES.Y/2, 48, rl.RayWhite)

			rl.EndDrawing()
		case utils.PAUSED:
			if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyM) {
				state.View = utils.IN_GAME
			}

			if rl.IsKeyPressed(rl.KeyQ) {
				gameState.AppState = nil
				state.View = utils.MAIN_MENU
			}

			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)
			rl.DrawText("Paused", state.RES.X/2, state.RES.Y/3, 48, rl.RayWhite)

			rl.EndDrawing()
		case utils.IN_GAME:
			game.GameUpdate(&state, &gameState, &character_textures, &tile_textures)
		}
	}

	for _, t := range tile_textures {
		rl.UnloadTexture(t)
	}
	for _, t := range character_textures {
		rl.UnloadTexture(t)
	}

	rl.CloseWindow()
}
