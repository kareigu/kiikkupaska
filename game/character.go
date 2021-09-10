package game

import (
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

func (player *Player) Draw() {
	texture := rendering.GetCharacterSprite(player.State)
	rl.DrawTexture(*texture, player.Pos.X, player.Pos.Y, rl.White)
}

type Enemy struct {
	Pos        utils.IVector2
	State      int
	LightLevel uint8
	Stats      Stats
	Turn       TurnData
}

func (enemy *Enemy) GetTurn() *TurnData {
	return &enemy.Turn
}

func (enemy *Enemy) GetStats() *Stats {
	return &enemy.Stats
}

func (enemy *Enemy) StartTurn() {
	enemy.Turn.Actions = 3
	enemy.Turn.Movement = enemy.Stats.Movement
	enemy.Turn.Done = false
}

func (enemy *Enemy) Draw() {
	texture := rendering.GetCharacterSprite(enemy.State)
	colour := rl.White
	colour.A = enemy.LightLevel
	rl.DrawTexture(*texture, enemy.Pos.X, enemy.Pos.Y, colour)
}

func (enemy *Enemy) VisibleToPlayer() bool {
	player := state.Player
	visrange := int32(player.Stats.Visibility) * TILE_SIZE
	if enemy.Pos.X > player.Pos.X+(visrange) || enemy.Pos.X < player.Pos.X-(visrange) || enemy.Pos.Y > player.Pos.Y+(visrange) || enemy.Pos.Y < player.Pos.Y-(visrange) {
		return false
	} else {
		return true
	}
}

func (enemy *Enemy) DistanceToPlayer() float32 {
	enemy_vec := enemy.Pos.ToVec2()
	player_vec := state.Player.Pos.ToVec2()
	distance := rl.Vector2Distance(enemy_vec, player_vec) / float32(TILE_SIZE)

	return distance
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
