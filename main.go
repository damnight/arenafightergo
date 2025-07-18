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

package main

import (
	"arenafightergo/ecsgo"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Arena Fighter Go")
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	co, err := ecsgo.NewCoordinator()
	if err != nil {
		log.Fatal(err)
	}
	g, err := co.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
