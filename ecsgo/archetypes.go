package ecsgo

import "github.com/hajimehoshi/ebiten/v2"

type Creature struct {
	pos Position
	vel Velocity

	health    Health
	baseSpeed BaseSpeed

	sprites []*Sprite
}

type Renderable struct {
	slice *ComponentSlice[ArchetypeID]
}

type SpriteSheet struct {
	slice map[SpriteID][]*ebiten.Image
}

type Game struct {
	co            *Coordinator
	Width, Height int
	CurrentLevel  *Level

	CamX, CamY float64
	CamScale   float64
	CamScaleTo float64

	MousePanX, MousePanY int

	Offscreen *ebiten.Image
}

type Level struct {
	Width, Height int
	TileSize      int

	Map [LEVEL_W][LEVEL_H]SpriteID
}
