package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func calculateLightLevel(distance float32, visibilityStat uint8) uint8 {
	distance_alpha := distance / float32(visibilityStat)
	colour := rl.ColorAlpha(rl.White, distance_alpha)
	// Reverse alpha to make closer objects brighter instead of darker
	colour.A = uint8(math.Abs(float64(colour.A) - 255.0))
	return colour.A
}
