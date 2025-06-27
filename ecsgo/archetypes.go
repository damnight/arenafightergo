package ecsgo

import "github.com/hajimehoshi/ebiten/v2"

type Creature struct {
	pos Position
	vel Velocity

	health    Health
	baseSpeed BaseSpeed

	sprites []*Sprite
}

type Orc struct {
	creature Creature
	color    Color
}

type Game struct {
	w, h         int
	currentLevel *Level

	camX, camY float64
	camScale   float64
	camScaleTo float64

	mousePanX, mousePanY int

	offscreen *ebiten.Image
}

type Level struct {
	width, height int

	tiles    [][]*Tile // (Y,X) array of tiles
	tileSize int
}
