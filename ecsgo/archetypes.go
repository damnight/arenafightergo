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
	slice map[SpriteID][]*Sprite
}

type Tile struct {
	pos      *Position
	sprites  *[]*Sprite
	tileType SpriteID
}

//func (t *Tile) Reset() {
//	t.pos = Position{}
//	t.sprites = []*Sprite{}
//	t.tileType = Default
//}
//
//func (t *Tile) AddSprites(s []*Sprite) {
//	for _, sp := range s {
//		t.sprites = append(t.sprites, sp)
//	}
//}
//
//func (t *Tile) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
//	for _, s := range t.sprites {
//		screen.DrawImage(s.img, options)
//	}
//}

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
