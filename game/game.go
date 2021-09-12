package game

import (
	"math"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const TILE_SIZE int32 = 32
const PLAYER_OFFSET_X int32 = 0
const PLAYER_OFFSET_Y int32 = 0

type GameState struct {
	AppState      *utils.State
	Camera        *rl.Camera2D
	Player        *Player
	Map           [][]*Tile
	Enemies       []*Enemy
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
			TileLightFx:     true,
		},
		SelectionMode: SelectionMode{
			Using: false,
			Pos:   player.Pos,
		},
		tempTimeSinceTurn: 0.0,
	}
	state.Map, state.Enemies = GenerateLevel()
	return &state
}

func initPlayerAndCam(state *utils.State) (*Player, *rl.Camera2D) {
	cam := rl.Camera2D{
		Offset: rl.Vector2{
			X: float32(state.Settings.Resolution.X / 2),
			Y: float32(state.Settings.Resolution.Y / 2),
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

		if state.Player.Turn.Done {
			for _, enemy := range state.Enemies {
				if !enemy.Turn.Done {
					enemy.DoAction()
				}
			}
		}

		var enemiesToDraw []*Enemy
		for i, enemy := range state.Enemies {
			if enemy.Health <= 0.0 {
				length := len(state.Enemies)
				if length > 0 {
					state.Enemies[i] = state.Enemies[length-1]
					state.Enemies = state.Enemies[:length-1]
				}
			}
			if enemy.VisibleToPlayer() {
				enemy.LightLevel = calculateLightLevel(enemy.DistanceToPlayer(), state.Player.Stats.Visibility)
				enemiesToDraw = append(enemiesToDraw, enemy)
			}
		}

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
					if tile.VisibleToPlayer(&enemiesToDraw) {
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

		for _, enemy := range enemiesToDraw {
			enemy.Draw()
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

func reverseRange(v float32) float32 {
	return float32(math.Abs((float64(v) - 1.0) / (0.0 - 1.0)))
}
