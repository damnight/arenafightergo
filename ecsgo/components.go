package ecsgo

// this is a data class
import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteID uint
type ComponentType uint

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
	VoidType ComponentType = iota
	PositionType
	SpriteType
)

type IComponent interface {
	Type() ComponentType
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

func (s Sprite) Type() ComponentType {
	return SpriteType
}

type Position struct {
	x, y, z float64
}

func (p Position) Type() ComponentType {
	return PositionType
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
