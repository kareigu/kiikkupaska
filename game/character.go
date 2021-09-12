package game

import (
	"log"
	"math"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	Draw()
}

func InVisRange(source utils.IVector2, target utils.IVector2, visStat uint8) bool {
	visrange := int32(visStat) * TILE_SIZE
	return !(target.X > source.X+(visrange) || target.X < source.X-(visrange) || target.Y > source.Y+(visrange) || target.Y < source.Y-(visrange))
}

type Player struct {
	Pos   utils.IVector2
	State int
	Stats Stats
	Turn  TurnData
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

func (player *Player) EndTurn() {
	for _, enemy := range state.Enemies {
		enemy.StartTurn()
	}
	player.Turn.Done = true
}

func (player *Player) Draw() {
	texture := rendering.GetCharacterSprite(player.State)
	rl.DrawTexture(*texture, player.Pos.X, player.Pos.Y, rl.White)
}

func (player *Player) Move() {

	if player.Turn.Movement > 0 {
		p_x := player.Pos.X
		p_y := player.Pos.Y

		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			p_x -= TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			p_x += TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
			p_y -= TILE_SIZE
		}
		if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
			p_y += TILE_SIZE
		}

		npos := utils.IVector2{X: p_x - PLAYER_OFFSET_X, Y: p_y - PLAYER_OFFSET_Y}

		if tile, ok := getTile(npos); ok {
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

func (player *Player) Attack(enemy *Enemy) {
	dmg := (float32(player.Stats.Strength) + float32(player.Stats.Dexterity)) * 1.2
	enemy.Health -= dmg
	player.Turn.Actions--
	log.Printf("Attacked enemy with %.2f damage, leaving %.2f health", dmg, enemy.Health)
}

type Enemy struct {
	Pos                utils.IVector2
	State              int
	Health             float32
	LightLevel         uint8
	LastKnownPlayerPos utils.IVector2
	Stats              Stats
	Turn               TurnData
}

func (enemy *Enemy) GetTurn() *TurnData {
	return &enemy.Turn
}

func (enemy *Enemy) GetStats() *Stats {
	return &enemy.Stats
}

func (enemy *Enemy) StartTurn() {
	enemy.Turn.Actions = 1
	enemy.Turn.Movement = enemy.Stats.Movement
	enemy.Turn.Done = false
}

func (enemy *Enemy) Draw() {
	texture := rendering.GetCharacterSprite(enemy.State)
	rl.DrawTexture(*texture, enemy.Pos.X, enemy.Pos.Y, rl.White)
}

func (enemy *Enemy) DoAction() {
	if enemy.Turn.Movement > 0 {
		enemy.Move()
	} else {
		enemy.Turn.Done = true
	}
}

func (enemy *Enemy) Move() {
	e_x := enemy.Pos.X
	e_y := enemy.Pos.Y

	/* 	if enemy.CanSeePlayer() {
	   		enemy.LastKnownPlayerPos = state.Player.Pos
	   	}

	   	if enemy.Pos == enemy.LastKnownPlayerPos {
	   		switch rl.GetRandomValue(0, 3) {
	   		case 0:
	   			e_x -= TILE_SIZE
	   		case 1:
	   			e_x += TILE_SIZE
	   		case 2:
	   			e_y -= TILE_SIZE
	   		case 3:
	   			e_y += TILE_SIZE
	   		}

	   	} else {
	   		enemy_pos := enemy.Pos
	   		diff := rl.Vector2Subtract(enemy_pos.ToVec2(), enemy.LastKnownPlayerPos.ToVec2())

	   		if diff.X != 0.0 {
	   			if math.Signbit(float64(diff.X)) {
	   				e_x -= TILE_SIZE
	   			} else {
	   				e_x += TILE_SIZE
	   			}
	   		} else if diff.Y != 0.0 {
	   			if math.Signbit(float64(diff.Y)) {
	   				e_y -= TILE_SIZE
	   			} else {
	   				e_y += TILE_SIZE
	   			}
	   		}

	   		log.Printf("Trying to move to x: %d y: %d", e_x/TILE_SIZE, e_y/TILE_SIZE)

	   	} */

	if enemy.CanSeePlayer() {
		enemy_pos := enemy.Pos
		diff := rl.Vector2Subtract(enemy_pos.ToVec2(), enemy.LastKnownPlayerPos.ToVec2())

		if diff.X != 0.0 {
			if math.Signbit(float64(diff.X)) {
				e_x -= TILE_SIZE
			} else {
				e_x += TILE_SIZE
			}
		} else if diff.Y != 0.0 {
			if math.Signbit(float64(diff.Y)) {
				e_y -= TILE_SIZE
			} else {
				e_y += TILE_SIZE
			}
		}

		log.Printf("Trying to move to x: %d y: %d", e_x/TILE_SIZE, e_y/TILE_SIZE)
	} else {
		if enemy.Pos == enemy.LastKnownPlayerPos {
			switch rl.GetRandomValue(0, 3) {
			case 0:
				e_x -= TILE_SIZE
			case 1:
				e_x += TILE_SIZE
			case 2:
				e_y -= TILE_SIZE
			case 3:
				e_y += TILE_SIZE
			}
		}
	}

	npos := utils.IVector2{X: e_x, Y: e_y}

	if tile, ok := getTile(npos); ok {
		if tile.Block {
			return
		}
	} else {
		return
	}

	if state.Player.Pos == npos {
		return
	}

	/* if enemy.LastKnownPlayerPos == enemy.Pos {
		enemy.LastKnownPlayerPos = npos
	} */

	enemy.Pos = npos
	enemy.Turn.Movement--
}

func (enemy *Enemy) VisibleToPlayer() bool {
	return InVisRange(state.Player.Pos, enemy.Pos, state.Player.Stats.Visibility)
}

func (enemy *Enemy) DistanceToPlayer() float32 {
	enemy_vec := enemy.Pos.ToVec2()
	player_vec := state.Player.Pos.ToVec2()
	distance := rl.Vector2Distance(enemy_vec, player_vec) / float32(TILE_SIZE)

	return distance
}

func (enemy *Enemy) CanSeePlayer() bool {
	return InVisRange(enemy.Pos, state.Player.Pos, enemy.Stats.Visibility)
}

func (enemy *Enemy) LightEmittedToTile(tile *Tile) uint8 {
	distance := tile.DistanceToEnemy(enemy)
	return calculateLightLevel(distance, enemy.Stats.Visibility)
}

func CreateRandomEnemy(pos utils.IVector2) *Enemy {
	stats := DefaultGoblinStats()
	new_enemy := Enemy{
		Pos:                pos,
		LastKnownPlayerPos: pos,
		Health:             float32(stats.Vitality) * 2.63,
		State:              rendering.GOBLIN_IDLE,
		Stats:              stats,
		Turn:               DefaultEnemyTurn(),
	}
	return &new_enemy
}

func DefaultEnemyTurn() TurnData {
	return TurnData{
		Movement: 0,
		Actions:  0,
		Done:     true,
	}
}

func DefaultGoblinStats() Stats {
	return Stats{
		Movement:   4,
		Visibility: 4,
		Vitality:   5,
		Strength:   3,
		Dexterity:  5,
	}
}
