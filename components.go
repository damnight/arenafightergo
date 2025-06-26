package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	Wall TileType = iota
	Statue
	Crown
	Floor
	Tube
	Portal
)

type Tile struct {
	pos Position

	sprite   []*Sprite
	tileType *TileType
}

type Position struct {
	x, y, z float64
}

type Sprite struct {
	img *ebiten.Image
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
