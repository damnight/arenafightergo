package ecsgo

import (
	"arenafightergo/sprites"
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
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

func LoadSpriteSheet(tileSize int) (*SpriteSheet, error) {

	img, _, err := image.Decode(bytes.NewReader(sprites.Spritesheet_png))
	if err != nil {
		return &SpriteSheet{}, err
	}

	sheet := ebiten.NewImageFromImage(img)

	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y+1)*tileSize, (x+1)*tileSize, y*tileSize)).(*ebiten.Image)
	}

	s := &SpriteSheet{
		slice: make(map[SpriteID][]*ebiten.Image),
	}

	s.slice[Default] = append(s.slice[Default], spriteAt(10, 7))
	s.slice[LightRock] = append(s.slice[LightRock], spriteAt(9, 0))
	s.slice[CrackedEarth1] = append(s.slice[CrackedEarth1], spriteAt(6, 1))
	s.slice[CrackedEarth2] = append(s.slice[CrackedEarth2], spriteAt(7, 1))
	s.slice[CrackedEarthWeeds1] = append(s.slice[CrackedEarthWeeds1], spriteAt(8, 1))
	s.slice[CrackedEarthWeeds2] = append(s.slice[CrackedEarthWeeds2], spriteAt(9, 1))
	s.slice[LightGrass] = append(s.slice[LightGrass], spriteAt(0, 2))
	s.slice[DarkGrass] = append(s.slice[DarkGrass], spriteAt(7, 3))
	s.slice[WaterLight1] = append(s.slice[WaterLight1], spriteAt(0, 10))
	s.slice[RockWall1] = append(s.slice[RockWall1], spriteAt(6, 5))
	s.slice[RockPeak1] = append(s.slice[RockPeak1], spriteAt(9, 5))

	return s, nil

}

// LoadSpriteSheet loads the embedded SpriteSheet.
//unc LoadSpriteSheet(TileSize int) (*SpriteSheet, error) {
//       img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
//       if err != nil {
//       	return &SpriteSheet{}, err
//       }
//
//       sheet := ebiten.NewImageFromImage(img)
//
//       // spriteAt returns a sprite at the provided coordinates.
//       spriteAt := func(x, y int) *ebiten.Image {
//       	return sheet.SubImage(image.Rect(x*TileSize, (y+1)*TileSize, (x+1)*TileSize, y*TileSize)).(*ebiten.Image)
//       }
//
//       // Populate SpriteSheet.
//       s := &SpriteSheet{
//       	slice: make(map[SpriteID][]*ebiten.Image),
//       }
//
//       s.slice[Floor] = append(s.slice[Floor], spriteAt(10, 4))
//       s.slice[Wall] = append(s.slice[Wall], spriteAt(2, 3))
//       s.slice[Statue] = append(s.slice[Statue], spriteAt(5, 4))
//       s.slice[Tube] = append(s.slice[Tube], spriteAt(3, 4))
//       s.slice[Crown] = append(s.slice[Crown], spriteAt(8, 6))
//       s.slice[Portal] = append(s.slice[Portal], spriteAt(5, 6))
//       s.slice[Knight] = append(s.slice[Knight], spriteAt(4, 7))
//
//       return s, nil
//

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.Width, l.Height
}

// NewLevel returns a new randomly generated Level.
//unc (co *Coordinator) NewLevelIsoTest() (*Level, error) {
//       // Create a 108x108 Level.
//       l := &Level{
//       	Width:    128,
//       	Height:   128,
//       	TileSize: 64,
//       }
//
//       z := 0.0
//
//       for y := 0; y < l.Height; y++ {
//       	for x := 0; x < l.Width; x++ {
//       		isBorderSpace := x == 0 || y == 0 || x == l.Width-1 || y == l.Height-1
//       		val := rand.Intn(1000)
//       		switch {
//       		case isBorderSpace || val < 275:
//       			co.CreateTile(float64(x), float64(y), z, Wall)
//       		case val < 285:
//       			co.CreateTile(float64(x), float64(y), z, Statue)
//       		case val < 288:
//       			co.CreateTile(float64(x), float64(y), z, Crown)
//       		case val < 289:
//       			co.CreateTile(float64(x), float64(y), z, Floor)
//       			co.CreateTile(float64(x), float64(y), z, Tube)
//       		case val < 290:
//       			co.CreateTile(float64(x), float64(y), z, Portal)
//       		default:
//       			co.CreateTile(float64(x), float64(y), z, Floor)
//       		}
//       	}
//       }
//
//       return l, nil
//

func (co *Coordinator) NewLevel1() (*Level, error) {
	levelMap, err := LoadMap(LEVEL1_MAP_PATH)
	if err != nil {
		fmt.Errorf("can't load level")
		return &Level{}, err
	}

	for x, row := range levelMap.Map {
		for y, tile := range row {
			co.CreateTile(float64(x), float64(y), 0, tile)
		}
	}

	return levelMap, nil

}

type LevelRead struct {
	Map [][]SpriteID
}

func LoadMap(path string) (*Level, error) {

	csvFile, err := os.Open(path)
	if err != nil {
		return &Level{}, err
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	readMap, err := r.ReadAll()
	if err != nil {
		return &Level{}, err
	}

	level := &Level{
		Width:    LEVEL_W,
		Height:   LEVEL_H,
		TileSize: TILE_SIZE,
		Map:      [LEVEL_W][LEVEL_H]SpriteID{},
	}

	rowo := (LEVEL_W - len(readMap[0])) / 2
	colo := (LEVEL_H - len(readMap)) / 2

	for x := 0; x < len(readMap[0]); x++ {
		for y := 0; y < len(readMap); y++ {
			n, err := strconv.Atoi(readMap[x][y])
			if err != nil {
				level.Map[x][y] = Default
			}
			level.Map[x+rowo][y+colo] = SpriteID(n)
		}
	}

	return level, nil

}
