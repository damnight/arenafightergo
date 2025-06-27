package ecsgo

// this is a data class
import "github.com/hajimehoshi/ebiten/v2"

type TileType uint

const (
	Wall TileType = iota
	Statue
	Crown
	Floor
	Tube
	Portal
)

type IComponent interface {
	CreateComponent()
}

type Tile struct {
	pos Position

	sprite   []*Sprite
	tileType *TileType
}

func (t *Tile) CreateComponent(pos Position, sprite []*Sprite, tileType *TileType) {

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
