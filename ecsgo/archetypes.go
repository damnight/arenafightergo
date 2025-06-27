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
}

type SpriteSheet struct {
	terrain  map[SpriteID][]*Sprite
	creature map[SpriteID][]*Sprite
}

type Game struct {
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
	Tiles         [][]*Tile
	TileSize      int
}
