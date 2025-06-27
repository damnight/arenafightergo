package ecsgo

import "sync"

type World struct {
	creatures *ComponentSlice[Creature]
}

type ComponentManager struct {
	tilePool sync.Pool
}

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		tilePool: sync.Pool{
			New: func() any {
				// TODO: maybe put the sprite collocator func (don´t load each time) inside tile here?
				return &Tile{}
			},
		},
	}
}

type ComponentSlice[T any] struct {
	slice []T
}
