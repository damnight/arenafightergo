package ecsgo

// this is a data class
import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteID uint16
type ComponentTypeID uint16

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

const (
	VoidType ComponentTypeID = iota
	PositionType
	SpriteType
	MAX_COMPONENTTYPE_ID
)

type IComponent interface {
	Type() ComponentTypeID
}

type IPooled interface {
	Reset()
	Get()
	Put()
}

type IRenderable interface {
}

type Sprite struct {
	img      []*ebiten.Image
	spriteID SpriteID
}

func (s Sprite) Type() ComponentTypeID {
	return SpriteType
}

type Position struct {
	x, y, z float64
}

func (p Position) Type() ComponentTypeID {
	return PositionType
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
