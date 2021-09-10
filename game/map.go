package game

import (
	"log"
	"math/rand"
	"rendering"
	"strings"
	"time"
	"utils"

	simplex "github.com/ojrac/opensimplex-go"
)

const ENEMY_SPAWN_RATE = 0.7

func GenerateLevel() ([][]*Tile, []*Enemy) {
	t := time.Now()
	tiles := generateTiles()
	enemies := placeEnemies(tiles)

	log.Println("Level generated in ", time.Since(t))
	return tiles, enemies
}

func placeEnemies(tiles [][]*Tile) []*Enemy {
	t := time.Now()
	var enemies []*Enemy
	for _, row := range tiles {
		for _, tile := range row {
			if tile != nil && tile.Type == rendering.TILE_FLOOR_SPAWN && tile.Pos != state.Player.Pos {
				if rand.Float32() < ENEMY_SPAWN_RATE {
					stats := DefaultGoblinStats()
					new_enemy := Enemy{
						Pos:    tile.Pos,
						Health: float32(stats.Vitality) * 2.63,
						State:  rendering.GOBLIN_IDLE,
						Stats:  stats,
						Turn:   DefaultEnemyTurn(),
					}
					enemies = append(enemies, &new_enemy)
				}
			}
		}
	}
	log.Println("Enemy placement complete in ", time.Since(t))
	return enemies
}

func generateTiles() [][]*Tile {
	const tileArrDimensions = 1000
	mapstring := generateMapString()
	t := time.Now()
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
			tile := charToTile(char, pos)
			tiles[x][y] = &tile
		}
	}
	log.Println("Tiles generated in ", time.Since(t))
	utils.DebugPrint(mapstring)
	return tiles
}

func generateMapString() string {
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

func checkTileVisibility(player *Player, tile *Tile) (uint8, bool) {
	visrange := int32(player.Stats.Visibility) * TILE_SIZE
	if tile.Pos.X > player.Pos.X+(visrange) || tile.Pos.X < player.Pos.X-(visrange) || tile.Pos.Y > player.Pos.Y+(visrange) || tile.Pos.Y < player.Pos.Y-(visrange) {
		return 0, false
	} else {
		distance := getTileDistanceToPlayer(player, tile)
		return calculateLightLevel(distance, player.Stats.Visibility), true
	}
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
