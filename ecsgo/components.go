package ecsgo

// this is a data class
import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TileType uint

const (
	Default TileType = iota
	Statue
	Crown
	Floor
	Tube
	Portal
	Wall
)

type IComponent interface {
}

type IPooled interface {
	Reset()
	Get()
	Put()
}

type Tile struct {
	pos      Position
	sprites  []*Sprite
	tileType TileType
}

func (t *Tile) Reset() {
	t.pos = Position{}
	t.sprites = []*Sprite{}
	t.tileType = Default
}

func (t *Tile) AddSprite(s *Sprite) {
	t.sprites = append(t.sprites, s)
}

func (t *Tile) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	for _, s := range t.sprites {
		screen.DrawImage(s.img, options)
	}
}

type Position struct {
	x, y, z float64
}

type Velocity struct {
	dx, dy float64
}

type Color struct {
	r, g, b, a uint8
}

type Health struct {
	health int
}

type BaseSpeed struct {
	speed int
}

type Sprite struct {
	img *ebiten.Image
}
