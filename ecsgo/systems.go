package ecsgo

import (
	"bytes"
	"image"
	_ "image/png"
	"math/rand"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

type Systems struct {
}

type EntityFactory struct {
}

//func (ef *EntityFactory)creatureCreator(cID CreatureID, args []T) []EntityID {
//
//	switch cID {
//		case OrcID:
//			return ef.createOrc(args)
//		default:
//			fmt.Print("default behaviour creatureCreator")
//
//	}
//
//}
//
//func (ef *EntityFactory)createOrc(args []T) []EntityID {
//
//
//
//
//}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(TileSize int) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*TileSize, (y+1)*TileSize, (x+1)*TileSize, y*TileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{
		slice: make(map[SpriteID][]*Sprite),
	}

	s.slice[Floor] = append(s.slice[Wall], &Sprite{img: spriteAt(10, 4)})
	s.slice[Wall] = append(s.slice[Statue], &Sprite{img: spriteAt(2, 3)})
	s.slice[Statue] = append(s.slice[Floor], &Sprite{img: spriteAt(5, 4)})
	s.slice[Tube] = append(s.slice[Tube], &Sprite{img: spriteAt(3, 4)})
	s.slice[Crown] = append(s.slice[Crown], &Sprite{img: spriteAt(8, 6)})
	s.slice[Portal] = append(s.slice[Portal], &Sprite{img: spriteAt(5, 6)})
	s.slice[Knight] = append(s.slice[Knight], &Sprite{img: spriteAt(4, 7)})

	return s, nil
}

//// Tile returns the tile at the provided coordinates, or nil.
//func (l *Level) Tile(x, y int) *Tile {
//	if x >= 0 && y >= 0 && x < l.Width && y < l.Height {
//		return l.Tiles[y][x]
//	}
//	return nil
//}

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.Width, l.Height
}

// NewLevel returns a new randomly generated Level.
func NewLevel() (*Level, error) {
	// Create a 108x108 Level.
	l := &Level{
		Width:    128,
		Height:   128,
		TileSize: 64,
	}

	// Load embedded SpriteSheet.
	sheet, err := LoadSpriteSheet(l.TileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded spritesheet: %s", err)
	}

	// Fill each tile with one or more sprites randomly.
	// TODO: create tiles in function and add to ecs
	l.Tiles = make([][]*Tile, l.Height)
	for y := 0; y < l.Height; y++ {
		l.Tiles[y] = make([]*Tile, l.Width)
		for x := 0; x < l.Width; x++ {
			t := &Tile{}
			isBorderSpace := x == 0 || y == 0 || x == l.Width-1 || y == l.Height-1
			val := rand.Intn(1000)
			switch {
			case isBorderSpace || val < 275:
				t.AddSprites(sheet.slice[Wall])
			case val < 285:
				t.AddSprites(sheet.slice[Statue])
			case val < 288:
				t.AddSprites(sheet.slice[Crown])
			case val < 289:
				t.AddSprites(sheet.slice[Floor])
				t.AddSprites(sheet.slice[Tube])
			case val < 290:
				t.AddSprites(sheet.slice[Portal])
			default:
				t.AddSprites(sheet.slice[Floor])
			}
			l.Tiles[y][x] = t
		}
	}

	return l, nil
}
