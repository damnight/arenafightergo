package ecsgo

import (
	"bytes"
	"image"
	_ "image/png"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

type SystemID uint32

const (
	SpriteLoader SystemID = iota
)

type SystemsManager struct {
	SystemsIndex []SystemID
}

func NewSystemsManager() (*SystemsManager, error) {
	sm := &SystemsManager{}
	return sm, nil
}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(TileSize int) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
	if err != nil {
		return &SpriteSheet{}, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*TileSize, (y+1)*TileSize, (x+1)*TileSize, y*TileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{
		slice: make(map[SpriteID][]*ebiten.Image),
	}

	s.slice[Floor] = append(s.slice[Floor], spriteAt(10, 4))
	s.slice[Wall] = append(s.slice[Wall], spriteAt(2, 3))
	s.slice[Statue] = append(s.slice[Statue], spriteAt(5, 4))
	s.slice[Tube] = append(s.slice[Tube], spriteAt(3, 4))
	s.slice[Crown] = append(s.slice[Crown], spriteAt(8, 6))
	s.slice[Portal] = append(s.slice[Portal], spriteAt(5, 6))
	s.slice[Knight] = append(s.slice[Knight], spriteAt(4, 7))

	return s, nil
}

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.Width, l.Height
}

// NewLevel returns a new randomly generated Level.
func (cp *ComponentManager) NewLevel() (*Level, error) {
	// Create a 108x108 Level.
	l := &Level{
		Width:    128,
		Height:   128,
		TileSize: 64,
	}

	z := 0.0

	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			isBorderSpace := x == 0 || y == 0 || x == l.Width-1 || y == l.Height-1
			val := rand.Intn(1000)
			switch {
			case isBorderSpace || val < 275:
				cp.CreateTile(float64(x), float64(y), z, Wall)
			case val < 285:
				cp.CreateTile(float64(x), float64(y), z, Statue)
			case val < 288:
				cp.CreateTile(float64(x), float64(y), z, Crown)
			case val < 289:
				cp.CreateTile(float64(x), float64(y), z, Floor)
				cp.CreateTile(float64(x), float64(y), z, Tube)
			case val < 290:
				cp.CreateTile(float64(x), float64(y), z, Portal)
			default:
				cp.CreateTile(float64(x), float64(y), z, Floor)
			}
		}
	}

	return l, nil
}
