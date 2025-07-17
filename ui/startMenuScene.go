package ui

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type ElemType int

const (
	Void ElemType = iota
	TitleBox
	Background
	Button
)

type ElemState int

const (
	Idle ElemState = iota
	Active
	Hover
	Pressed
	Released
)

type ElemID int

const (
	Background_1 ElemID = iota
	Button_1
	Title
)

func MakeSceneManager() (*stagehand.SceneManager[StartMenuState], error) {
	s := &StartMenuScene{elems: make(map[ElemID]*ebiten.Image)}
	state := StartMenuState{}
	sm := stagehand.NewSceneManager[StartMenuState](s, state)

	return sm, nil

}

type StartMenuState struct {
	// state
}

type StartMenuScene struct {
	// scene
	state StartMenuState
	elems map[ElemID]*ebiten.Image
	sm    *stagehand.SceneManager[StartMenuState]
}

func (s *StartMenuScene) Update() error {
	// update scene
	return nil
}

func (s *StartMenuScene) Draw(screen *ebiten.Image) {
	// draw scene
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.elems[Background_1], op)
}

func (s *StartMenuScene) Load(state StartMenuState, sm stagehand.SceneController[StartMenuState]) {
	// load scene
	img, _, err := image.Decode(bytes.NewReader(Background_png))
	if err != nil {
		panic(err)
	}

	bg := ebiten.NewImageFromImage(img)
	s.elems[Background_1] = bg

	s.sm = sm.(*stagehand.SceneManager[StartMenuState])

}

func (s *StartMenuScene) Unload() StartMenuState {
	// unload
	return s.state

}
func (s *StartMenuScene) Layout(w, h int) (int, int) {
	return w, h
}
