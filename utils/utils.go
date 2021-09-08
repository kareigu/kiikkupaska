package utils

const (
	MAIN_MENU = iota
	PAUSED    = iota
	IN_GAME   = iota
)

type State struct {
	Loading bool
	View    int
	RES     IVector2
}

type IVector2 struct {
	X int32
	Y int32
}
