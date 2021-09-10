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
	texture := rendering.GetPlayerSprite(player.State)
	rl.DrawTexture(*texture, player.Pos.X, player.Pos.Y, rl.White)
}

type Enemy struct {
	Pos   utils.IVector2
	State int
	Stats Stats
	Turn  TurnData
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
