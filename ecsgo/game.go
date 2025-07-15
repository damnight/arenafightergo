// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecsgo

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (co *Coordinator) NewGame() (*Game, error) {
	l, err := co.NewLevel()
	if err != nil {
		fmt.Println("NEW GAME ERROR")
		return nil, fmt.Errorf("failed to create new level: %s", err)
	}

	g := &Game{
		co:           co,
		CurrentLevel: l,
		CamScale:     1,
		CamScaleTo:   1,
		MousePanX:    math.MinInt32,
		MousePanY:    math.MinInt32,
	}
	return g, nil
}

// Update reads current user input and updates the Game state.
func (g *Game) Update() error {
	// Update target zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = .25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	g.CamScaleTo += scrollY * (g.CamScaleTo / 7)

	// Clamp target zoom level.
	if g.CamScaleTo < 0.01 {
		g.CamScaleTo = 0.01
	} else if g.CamScaleTo > 100 {
		g.CamScaleTo = 100
	}

	// Smooth zoom transition.
	div := 10.0
	if g.CamScaleTo > g.CamScale {
		g.CamScale += (g.CamScaleTo - g.CamScale) / div
	} else if g.CamScaleTo < g.CamScale {
		g.CamScale -= (g.CamScale - g.CamScaleTo) / div
	}

	// Pan camera via keyboard.
	pan := 7.0 / g.CamScale
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.CamX -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.CamX += pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.CamY -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.CamY += pan
	}

	// Pan camera via mouse.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if g.MousePanX == math.MinInt32 && g.MousePanY == math.MinInt32 {
			g.MousePanX, g.MousePanY = ebiten.CursorPosition()
		} else {
			x, y := ebiten.CursorPosition()
			dx, dy := float64(g.MousePanX-x)*(pan/100), float64(g.MousePanY-y)*(pan/100)
			g.CamX, g.CamY = g.CamX-dx, g.CamY+dy
		}
	} else if g.MousePanX != math.MinInt32 || g.MousePanY != math.MinInt32 {
		g.MousePanX, g.MousePanY = math.MinInt32, math.MinInt32
	}

	// Clamp camera position.
	worldWidth := float64(g.CurrentLevel.Width * g.CurrentLevel.TileSize / 2)
	worldHeight := float64(g.CurrentLevel.Height * g.CurrentLevel.TileSize / 2)
	if g.CamX < -worldWidth {
		g.CamX = -worldWidth
	} else if g.CamX > worldWidth {
		g.CamX = worldWidth
	}
	if g.CamY < -worldHeight {
		g.CamY = -worldHeight
	} else if g.CamY > 0 {
		g.CamY = 0
	}

	// Randomize level.
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		l, err := g.co.NewLevel()
		if err != nil {
			return fmt.Errorf("failed to create new level: %s", err)
		}

		g.CurrentLevel = l
	}

	return nil
}

// Draw draws the Game on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Render level.

	start := time.Now()
	g.co.renderLevel(screen, g)
	finish := time.Since(start)
	fmt.Printf("| Render Level: %v | Drawcalls: %v |\n", finish, drawcalls)
	drawcalls = 0
	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC R\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.CamScale, g.CamX, g.CamY))
}

// Layout is called when the Game's layout changes.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.Width, g.Height = outsideWidth, outsideHeight
	return g.Width, g.Height
}

// cartesianToIso transforms cartesian coordinates into isometric coordinates.
func (g *Game) cartesianToIso(x, y float64) (float64, float64) {
	TileSize := g.CurrentLevel.TileSize
	ix := (x - y) * float64(TileSize/2)
	iy := (x + y) * float64(TileSize/4)
	return ix, iy
}

/*
This function might be useful for those who want to modify this example.

// isoToCartesian transforms isometric coordinates into cartesian coordinates.
func (g *Game) isoToCartesian(x, y float64) (float64, float64) {
	TileSize := g.CurrentLevel.TileSize
	cx := (x/float64(TileSize/2) + y/float64(TileSize/4)) / 2
	cy := (y/float64(TileSize/4) - (x / float64(TileSize/2))) / 2
	return cx, cy
}
*/
