package ecsgo

import (
	"fmt"
	"slices"
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
	ComponentDefinitions map[*IComponent]ComponentID

	spriteSheet *SpriteSheet
}

func (cp *ComponentManager) GetTileSprites(tileType SpriteID) *[]*Sprite {

	tileSprites := cp.spriteSheet.slice[tileType]
	return &tileSprites
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
	cp := &ComponentManager{}

	sheet, err := LoadSpriteSheet(64)
	if err != nil {
		sheet := &SpriteSheet{}

		cp.spriteSheet = sheet
		return cp, fmt.Errorf("wasn´t able to load spritesheet: %s", err)
	}

	cp.spriteSheet = sheet

	return cp, nil

}
func (cp *ComponentManager) CheckForArchetype(comps []IComponent) ArchetypeID {

	// get IDs from ComponentDefinitions
	// get all associated Archetypes by ComponentIndex

	archList := []ArchetypeID{}
	compList := []ComponentID{}

	for _, c := range comps {
		id := cp.ComponentDefinitions[&c]
		archetypes := cp.ComponentIndex[id]
		archList = archetypes
	}

	equality := false
	for _, arch := range archList {
		archtypeComponents := cp.ArchetypeDefinitions[arch]
		equality = slices.Equal(archtypeComponents, compList)
		if equality {
			return arch
		}
	}
	// only executes if no match was found, thus creating a new archetype
	archID := CreateArchetypeID()
	cp.ArchetypeDefinitions[archID] = compList
	return archID
}

func (cp *ComponentManager) AddEntity(e EntityID, comps []IComponent) {
	
	// coordinate entity and components
	cp.EntityIndex[e] =


	// check if archetype exists for this set of Components
	arch := cp.CheckForArchetype(comps)



	AddToWorld()

}

type ComponentSlice[T any] struct {
	slice map[ArchetypeID][]T
}
