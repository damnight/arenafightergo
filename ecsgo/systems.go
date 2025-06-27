package ecsgo

//import "fmt"

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
func LoadSpriteSheet(tileSize int) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y+1)*tileSize, (x+1)*tileSize, y*tileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{}
	s.Floor = spriteAt(10, 4)
	s.Wall = spriteAt(2, 3)
	s.Statue = spriteAt(5, 4)
	s.Tube = spriteAt(3, 4)
	s.Crown = spriteAt(8, 6)
	s.Portal = spriteAt(5, 6)
	s.Knight = spriteAt(4, 7)

	return s, nil
}

// Tile returns the tile at the provided coordinates, or nil.
func (l *Level) Tile(x, y int) *Tile {
	if x >= 0 && y >= 0 && x < l.w && y < l.h {
		return l.tiles[y][x]
	}
	return nil
}

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.w, l.h
}

// NewLevel returns a new randomly generated Level.
func NewLevel() (*Level, error) {
	// Create a 108x108 Level.
	l := &Level{
		w:        128,
		h:        128,
		tileSize: 64,
	}

	// Load embedded SpriteSheet.
	ss, err := LoadSpriteSheet(l.tileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded spritesheet: %s", err)
	}

	// Fill each tile with one or more sprites randomly.
	l.tiles = make([][]*Tile, l.h)
	for y := 0; y < l.h; y++ {
		l.tiles[y] = make([]*Tile, l.w)
		for x := 0; x < l.w; x++ {
			t := &Tile{}
			isBorderSpace := x == 0 || y == 0 || x == l.w-1 || y == l.h-1
			val := rand.IntN(1000)
			switch {
			case isBorderSpace || val < 275:
				t.AddSprite(ss.Wall)
			case val < 285:
				t.AddSprite(ss.Statue)
			case val < 288:
				t.AddSprite(ss.Crown)
			case val < 289:
				t.AddSprite(ss.Floor)
				t.AddSprite(ss.Tube)
			case val < 290:
				t.AddSprite(ss.Portal)
			default:
				t.AddSprite(ss.Floor)
			}
			l.tiles[y][x] = t
		}
	}

	return l, nil
}
