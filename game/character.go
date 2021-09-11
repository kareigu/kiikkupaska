package game

import (
	"log"
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

func (player *Player) Attack(enemy *Enemy) {
	dmg := (float32(player.Stats.Strength) + float32(player.Stats.Dexterity)) * 1.2
	enemy.Health -= dmg
	player.Turn.Actions--
	log.Printf("Attacked enemy with %.2f damage, leaving %.2f health", dmg, enemy.Health)
}

type Enemy struct {
	Pos        utils.IVector2
	State      int
	Health     float32
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
	rl.DrawTexture(*texture, enemy.Pos.X, enemy.Pos.Y, rl.White)
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

func (enemy *Enemy) LightEmittedToTile(tile *Tile) uint8 {
	distance := tile.DistanceToEnemy(enemy)
	return calculateLightLevel(distance, enemy.Stats.Visibility)
}

func CreateRandomEnemy(pos utils.IVector2) *Enemy {
	stats := DefaultGoblinStats()
	new_enemy := Enemy{
		Pos:    pos,
		Health: float32(stats.Vitality) * 2.63,
		State:  rendering.GOBLIN_IDLE,
		Stats:  stats,
		Turn:   DefaultEnemyTurn(),
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
