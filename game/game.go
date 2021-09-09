package game

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
	simplex "github.com/ojrac/opensimplex-go"
)

const TILE_SIZE int32 = 32
const PLAYER_OFFSET_X int32 = 0
const PLAYER_OFFSET_Y int32 = 0

type GameState struct {
	AppState      *utils.State
	Camera        *rl.Camera2D
	Player        *Player
	Map           *map[utils.IVector2]Tile
	SelectionMode SelectionMode
	DebugDisplay  DebugDisplayData

	tempTimeSinceTurn float32
}

type SelectionMode struct {
	Using bool
	Pos   utils.IVector2
}

type Tile struct {
	Texture rl.Texture2D
	Pos     utils.IVector2
	Block   bool
}

type Player struct {
	Pos    utils.IVector2
	Sprite rl.Texture2D
	Stats  Stats
	Turn   TurnData
}

type TurnData struct {
	Movement uint8
	Actions  uint8
	Done     bool
}

type Stats struct {
	Movement   uint8
	Visibility uint8
	Vitality   uint8
	Strength   uint8
	Dexterity  uint8
}

type Character interface {
	GetTurn() *TurnData
	GetStats() *Stats
	StartTurn()
}

func (player *Player) GetTurn() *TurnData {
	return &player.Turn
}

func (player *Player) GetStats() *Stats {
	return &player.Stats
}

func (player *Player) StartTurn() {
	player.Turn.Actions = 3
	player.Turn.Movement = player.Stats.Movement
	player.Turn.Done = false
}

var state GameState

func InitGame(appState *utils.State, character_textures *[]rl.Texture2D, tile_textures *[]rl.Texture2D) *GameState {
	player, cam := initPlayerAndCam(appState, character_textures)
	state = GameState{
		AppState: appState,
		Player:   player,
		Camera:   cam,
		Map:      nil,
		DebugDisplay: DebugDisplayData{
			Enabled:         false,
			TileDisplayMode: DD_TILE_DISTANCE_FROM_PLAYER,
		},
		SelectionMode: SelectionMode{
			Using: false,
			Pos:   player.Pos,
		},
		tempTimeSinceTurn: 0.0,
	}
	state.Map = generateTiles(tile_textures)
	return &state
}

func initPlayerAndCam(state *utils.State, character_textures *[]rl.Texture2D) (*Player, *rl.Camera2D) {
	cam := rl.Camera2D{
		Offset: rl.Vector2{
			X: float32(state.RES.X / 2),
			Y: float32(state.RES.Y / 2),
		},
		Target: rl.Vector2{
			X: 0.0,
			Y: 0.0,
		},
		Rotation: 0.0,
		Zoom:     1.0,
	}

	player := Player{
		Pos:    utils.IVector2{X: PLAYER_OFFSET_X, Y: PLAYER_OFFSET_Y},
		Sprite: (*character_textures)[0],
		Stats: Stats{
			Movement:   6,
			Visibility: 8,
			Vitality:   6,
			Strength:   6,
			Dexterity:  6,
		},
		Turn: TurnData{
			Movement: 6,
			Actions:  3,
			Done:     false,
		},
	}

	return &player, &cam
}

func GameUpdate(appState *utils.State, gameState **GameState, character_textures *[]rl.Texture2D, tile_textures *[]rl.Texture2D, ui_sprites *[]rl.Texture2D) {
	if state.AppState == nil {
		appState.Loading = true
		*gameState = InitGame(appState, character_textures, tile_textures)
		appState.Loading = false
	} else {

		if !state.Player.Turn.Done {
			if state.SelectionMode.Using {
				moveSelectionCursor(&state.SelectionMode)
			} else {
				movePlayer(state.Player, state.Map)
			}
			state.tempTimeSinceTurn = 0.0
		} else {
			if state.tempTimeSinceTurn > 5.0 {
				state.Player.StartTurn()
			} else {
				state.tempTimeSinceTurn += rl.GetFrameTime()
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			state.SelectionMode.Pos = state.Player.Pos
			state.SelectionMode.Using = !state.SelectionMode.Using
		}

		if rl.IsKeyPressed(rl.KeyI) {
			state.Player.Stats.Visibility++
		}

		if rl.IsKeyPressed(rl.KeyK) {
			state.Player.Stats.Visibility--
		}

		if rl.IsKeyPressed(rl.KeyM) || rl.IsKeyPressed(rl.KeyEscape) {
			state.AppState.View = utils.PAUSED
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			state.DebugDisplay.Enabled = !state.DebugDisplay.Enabled
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			state.Player.Turn.Done = true
		}

		state.Camera.Target.X = float32(state.Player.Pos.X)
		state.Camera.Target.Y = float32(state.Player.Pos.Y)

		rl.BeginDrawing()

		//
		//	Draw 2D objects
		//	Characters, tiles etc.
		//
		rl.BeginMode2D(*state.Camera)
		rl.ClearBackground(rl.Black)

		for _, tile := range *state.Map {
			// Check if tile coordinates are in player visibility range
			// If not don't bother rendering it
			if colour, ok := checkTileVisibility(state.Player, &tile); ok {
				rl.DrawTexture(tile.Texture, tile.Pos.X, tile.Pos.Y, colour)

				if state.DebugDisplay.Enabled {
					handleTileDebugDisplay(&tile, colour)
				}
			}
		}

		rl.DrawTexture(state.Player.Sprite, state.Player.Pos.X, state.Player.Pos.Y, rl.White)
		if state.SelectionMode.Using {
			alpha := float32((math.Cos(3.0*float64(rl.GetTime())) + 1) * 0.5)
			rl.DrawTexture((*ui_sprites)[2], state.SelectionMode.Pos.X, state.SelectionMode.Pos.Y, rl.ColorAlpha(rl.White, alpha))
		}

		rl.EndMode2D()

		//
		//	Draw UI stuff
		//

		if !state.Player.Turn.Done {
			for h := 0; h < int(state.Player.Turn.Actions); h++ {
				rl.DrawTexture((*ui_sprites)[0], state.AppState.RES.X-100, int32(10+h*int(TILE_SIZE)+5), rl.White)
			}

			for m := 0; m < int(state.Player.Turn.Movement); m++ {
				rl.DrawTexture((*ui_sprites)[1], state.AppState.RES.X-60, int32(10+m*int(TILE_SIZE)+5), rl.White)
			}

			if !(state.Player.Turn.Actions > 0) || !(state.Player.Turn.Movement > 0) {
				utils.DrawMainText(rl.NewVector2(float32(state.AppState.RES.X)/2.7, float32(state.AppState.RES.Y)/1.1), 48.0, "ENTER TO END TURN", rl.RayWhite)
			}
		} else {
			utils.DrawMainText(rl.NewVector2(float32(state.AppState.RES.X)/2.7, float32(state.AppState.RES.Y)/8.0), 48.0, "PROCESSING TURNS", rl.RayWhite)
		}

		if state.DebugDisplay.Enabled {
			drawDebugSettings()
			drawDebugInfo()
		}
		rl.EndDrawing()
	}
}

func generateTiles(tile_textures *[]rl.Texture2D) *map[utils.IVector2]Tile {
	mapstring := generateMap()
	tiles := make(map[utils.IVector2]Tile)
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
			pos := utils.IVector2{X: pos_x, Y: pos_y}
			tile := charToTile(tile_textures, char, pos)
			tiles[pos] = tile
		}
	}
	utils.DebugPrint(mapstring)
	return &tiles
}

func checkTileVisibility(player *Player, tile *Tile) (rl.Color, bool) {
	visrange := int32(player.Stats.Visibility) * TILE_SIZE
	if tile.Pos.X > player.Pos.X+(visrange) || tile.Pos.X < player.Pos.X-(visrange) || tile.Pos.Y > player.Pos.Y+(visrange) || tile.Pos.Y < player.Pos.Y-(visrange) {
		return rl.Black, false
	} else {
		distance_alpha := float32(getTileDistanceToPlayer(player, tile)) / float32(player.Stats.Visibility)
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

func IVec2ToVec2(ivec utils.IVector2) rl.Vector2 {
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

func charToTile(texturelist *[]rl.Texture2D, c string, pos utils.IVector2) Tile {
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

func movePlayer(player *Player, tiles *map[utils.IVector2]Tile) {

	if player.Turn.Movement > 0 {
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

		npos := utils.IVector2{X: p_x - PLAYER_OFFSET_X, Y: p_y - PLAYER_OFFSET_Y}

		if tile, ok := (*tiles)[npos]; ok {
			if tile.Block {
				return
			}
		} else {
			return
		}

		if p_x != player.Pos.X || p_y != player.Pos.Y {
			player.Pos.X = p_x
			player.Pos.Y = p_y
			player.Turn.Movement--
		}
	}
}

func moveSelectionCursor(selection *SelectionMode) {
	s_x := selection.Pos.X
	s_y := selection.Pos.Y

	if rl.IsKeyPressed(rl.KeyLeft) {
		s_x -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		s_x += TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		s_y -= TILE_SIZE
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		s_y += TILE_SIZE
	}

	selection.Pos.X = s_x
	selection.Pos.Y = s_y
}
