package ecsgo

import (
	"fmt"
	"sync"
)

type World struct {
	positions *ComponentSlice[Position]
	tiles     *ComponentSlice[Tile]
}

type ComponentManager struct {
	// EntityID -> []ComponentID
	// ComponentID -> []ArchetypID
	// ArchetypeID -> []EntityID
	EntityIndex          map[EntityID][]ComponentID
	ComponentIndex       map[ComponentID][]ArchetypeID
	ArchetypeIndex       map[ArchetypeID][]EntityID
	ArchetypeDefinitions map[ArchetypeID][]ComponentID
	ComponentDefinitions map[ComponentID]*IComponent

	tilePool    sync.Pool
	spritePool  sync.Pool
	spriteSheet *SpriteSheet
}

func AddToWorld(comp IComponent, world World) {
	switch comp.(type) {
	case Position:
		world.positions = append(world.positions, comp)
	case Tile:
		world.tiles = append(world.tiles, comp)
	default:
	}

}

func NewComponentManager() (*ComponentManager, error) {
	cp := &ComponentManager{
		tilePool: sync.Pool{
			New: func() any {
				return &Tile{}
			},
		},
		spritePool: sync.Pool{
			New: func() any {
				return &Sprite{}
			},
		},
	}

	sheet, err := LoadSpriteSheet(64)
	if err != nil {
		sheet := &SpriteSheet{}

		cp.spriteSheet = sheet
		return cp, fmt.Errorf("wasn´t able to load spritesheet: %s", err)
	}

	cp.spriteSheet = sheet

	return cp, nil

}

func (cp *ComponentManager) AddEntity(e EntityID, arch ArchetypeID) {
	//
	e := CreateEntity()
	cID := CreateComponentID()

	cp.ComponentDefinitions[cID] = &comp
	cp.EntityIndex[e] = append(cp.EntityIndex[e], cID)

	cp.ComponentIndex[cID] = append(cp.ComponentIndex[cID], arch)
	cp.ArchetypeIndex[arch] = append(cp.ArchetypeIndex[arch], e)
	cp.ArchetypeDefinitions[arch] = append(cp.ArchetypeDefinitions[arch], cID)
	// add component to the world
	AddToWorld(comp)

}

type ComponentSlice[T any] struct {
	slice map[ArchetypeID][]T
}
