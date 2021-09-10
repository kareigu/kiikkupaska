package game

import (
	"log"
	"math"
	"math/rand"
	"rendering"
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
	Map           [][]*Tile
	SelectionMode SelectionMode
	DebugDisplay  DebugDisplayData

	tempTimeSinceTurn float32
}

type SelectionMode struct {
	Using bool
	Pos   utils.IVector2
}

var state GameState

func InitGame(appState *utils.State) *GameState {
	player, cam := initPlayerAndCam(appState)
	state = GameState{
		AppState: appState,
		Player:   player,
		Camera:   cam,
		Map:      nil,
		DebugDisplay: DebugDisplayData{
			Enabled:         false,
			TileDisplayMode: DD_TILE_NO_DISPLAY,
		},
		SelectionMode: SelectionMode{
			Using: false,
			Pos:   player.Pos,
		},
		tempTimeSinceTurn: 0.0,
	}
	state.Map = generateTiles(&appState.RenderAssets.TileTextures)
	return &state
}

func initPlayerAndCam(state *utils.State) (*Player, *rl.Camera2D) {
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
		Pos:   utils.IVector2{X: PLAYER_OFFSET_X, Y: PLAYER_OFFSET_Y},
		State: rendering.PLAYER_IDLE,
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

func GameUpdate(appState *utils.State, gameState **GameState) {
	if state.AppState == nil {
		appState.Loading = true
		*gameState = InitGame(appState)
		appState.Loading = false
	} else {
		HandleControls()

		//*
		//*	Filter out the tiles that are visible to the player
		//*	If the the tile is visible push it to a separate array
		//*	that the renderer can use to save time not going through all this at render time
		//*
		var tilesToDraw []*Tile

		for _, tile_row := range state.Map {
			for _, tile := range tile_row {
				if tile != nil {
					//! Check if tile coordinates are in the player's visibility range
					//! If not, don't bother adding it for render
					if lightLevel, ok := checkTileVisibility(state.Player, tile); ok {
						tile.LightLevel = lightLevel
						tilesToDraw = append(tilesToDraw, tile)
					}
				}
			}
		}

		rl.BeginDrawing()

		//*
		//*	Draw 2D objects
		//*	Characters, tiles etc.
		//*
		rl.BeginMode2D(*state.Camera)
		rl.ClearBackground(rl.Black)

		for _, tile := range tilesToDraw {
			tile.Draw()

			if state.DebugDisplay.Enabled {
				handleTileDebugDisplay(tile)
			}
		}

		state.Player.Draw()
		drawSelectionCursor()

		rl.EndMode2D()

		//*
		//*	UI Section
		//*
		drawUI()

		rl.EndDrawing()
	}
}

func generateTiles(tile_textures *[]rl.Texture2D) [][]*Tile {
	const tileArrDimensions = 1000
	mapstring := generateMap()
	player := state.Player
	tiles := make([][]*Tile, tileArrDimensions)
	for i := 0; i < tileArrDimensions; i++ {
		tiles[i] = make([]*Tile, tileArrDimensions)
	}

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
			tiles[x][y] = &tile
		}
	}
	utils.DebugPrint(mapstring)
	return tiles
}

func checkTileVisibility(player *Player, tile *Tile) (uint8, bool) {
	visrange := int32(player.Stats.Visibility) * TILE_SIZE
	if tile.Pos.X > player.Pos.X+(visrange) || tile.Pos.X < player.Pos.X-(visrange) || tile.Pos.Y > player.Pos.Y+(visrange) || tile.Pos.Y < player.Pos.Y-(visrange) {
		return 0, false
	} else {
		distance_alpha := float32(getTileDistanceToPlayer(player, tile)) / float32(player.Stats.Visibility)
		colour := rl.ColorAlpha(rl.White, distance_alpha)
		// Reverse alpha to make closer tiles brighter instead of darker
		colour.A = uint8(math.Abs(float64(colour.A) - 255.0))
		return colour.A, true
	}
}

func getTileDistanceToPlayer(player *Player, tile *Tile) float32 {
	tile_vec := tile.Pos.ToVec2()
	player_vec := player.Pos.ToVec2()
	distance := rl.Vector2Distance(tile_vec, player_vec) / float32(TILE_SIZE)

	return distance
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

func getTile(pos utils.IVector2) (*Tile, bool) {
	x := pos.X / TILE_SIZE
	y := pos.Y / TILE_SIZE

	if tile := state.Map[x][y]; tile != nil {
		return tile, true
	} else {
		return nil, false
	}
}
