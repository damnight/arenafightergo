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

type IRenderable interface {
}

type Sprite struct {
	img *ebiten.Image
}

type Position struct {
	x, y, z float64
}

func CreatePosition(x, y, z float64) (*Position, ComponentID) {
	pos := &Position{x, y, z}
	c := CreateComponentID()
	return pos, c
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
