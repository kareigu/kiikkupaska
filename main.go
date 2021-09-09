package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"game"
	"rendering"
	"utils"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var state utils.State

func init() {
	runtime.LockOSThread()
}

// Handle cli arguments
func init() {
	widthFlag := flag.Int("w", 800, "Define the window width")
	heightFlag := flag.Int("h", 600, "Define the window height")
	musicFlag := flag.Bool("music", false, "Enable or disable music")

	flag.Parse()

	log.Printf("Running with flags: -w %d -h %d -music=%v", *widthFlag, *heightFlag, *musicFlag)

	state = utils.State{
		Loading: false,
		View:    utils.MAIN_MENU,
		RES:     utils.IVector2{X: int32(*widthFlag), Y: int32(*heightFlag)},
		Music:   *musicFlag,
	}
}

var music = make([]rl.Music, 4)

func main() {
	rl.InitWindow(state.RES.X, state.RES.Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	tile_textures := rendering.LoadTileTextures()
	character_textures := rendering.LoadCharacterTextures()
	font := rendering.LoadFont()
	rendering.LoadGUIStylesheet()
	rl.InitAudioDevice()

	var gameState *game.GameState

	midx := 0
	if state.Music {
		music[0] = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "sng01_int01.wav"))
		music[1] = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "sng01_pre01.wav"))
		music[3] = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "sng01_cbt01.wav"))
		music[2] = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "sng01_aft01.wav"))
		music[0].Looping = false
		music[1].Looping = false
		music[2].Looping = false

		rl.PlayMusicStream(music[0])
		rl.PlayMusicStream(music[1])
		rl.PlayMusicStream(music[2])
		rl.PlayMusicStream(music[3])
	}

	for !rl.WindowShouldClose() {
		rl.SetWindowTitle(fmt.Sprintf("kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))

		if state.Music {
			rl.UpdateMusicStream(music[midx])
			if !rl.IsMusicStreamPlaying(music[midx]) {
				midx++
			}
		}

		switch state.View {
		case utils.MAIN_MENU:
			if rl.IsKeyPressed(rl.KeyEnter) {
				state.View = utils.IN_GAME
			}

			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)
			rl.DrawTextEx(font, "Main Menu", rl.Vector2{X: float32(state.RES.X/2 - 170), Y: float32(state.RES.Y / 6)}, 96.0, 1.0, rl.RayWhite)
			start := rgui.Button(rl.Rectangle{X: float32(state.RES.X) / 2.0, Y: float32(state.RES.Y) / 2.0, Width: 100.0, Height: 25.0}, "Start")
			exit := rgui.Button(rl.Rectangle{X: float32(state.RES.X) / 2.0, Y: float32(state.RES.Y)/2.0 + 100.0, Width: 100.0, Height: 25.0}, "Quit")

			rl.EndDrawing()
			if start {
				state.View = utils.IN_GAME
			}
			if exit {
				cleanup(&tile_textures, &character_textures)
				rl.CloseWindow()
			}
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
			rl.DrawTextEx(font, "Paused", rl.Vector2{X: float32(state.RES.X/2 - 170), Y: float32(state.RES.Y / 6)}, 96.0, 1.0, rl.RayWhite)

			rl.EndDrawing()
		case utils.IN_GAME:
			game.GameUpdate(&state, &gameState, &character_textures, &tile_textures)
		}
	}

	cleanup(&tile_textures, &character_textures)

	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func cleanup(tile_textures *[]rl.Texture2D, character_textures *[]rl.Texture2D) {
	for _, t := range *tile_textures {
		rl.UnloadTexture(t)
	}
	for _, t := range *character_textures {
		rl.UnloadTexture(t)
	}
}
