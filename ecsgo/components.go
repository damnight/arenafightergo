package ecsgo

// this is a data class
import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteID uint

const (
	Default SpriteID = iota
	Statue
	Crown
	Floor
	Tube
	Portal
	Wall
	Knight
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
	tileType SpriteID
}

func (t *Tile) Reset() {
	t.pos = Position{}
	t.sprites = []*Sprite{}
	t.tileType = Default
}

func (t *Tile) AddSprites(s []*Sprite) {
	for _, sp := range s {
		t.sprites = append(t.sprites, sp)
	}
}

func (t *Tile) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	for _, s := range t.sprites {
		screen.DrawImage(s.img, options)
	}
}

type Sprite struct {
	img *ebiten.Image
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
