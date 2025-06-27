package ecsgo

import (
	"fmt"
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

func (cp *ComponentManager) GetTileSprites(tileType SpriteID) *[]*Sprite{
	
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

	archList := make(map[ComponentID][]ArchetypeID)
	for _, c := range comps {
		id := cp.ComponentDefinitions[&c]
		archetypes := cp.ComponentIndex[id]
		archList[id] = archetypes
	}
	
	for , arch := range archList {

	}

	
	// make subset of Archetypes where all components are in
	


	
	// take smallest subset and return that archetype as ID as result


}


func (cp *ComponentManager) AddEntity(e EntityID, comps []IComponent) {
	
	// check if archetype exists for this set of Components
	CheckForArchetype()

	cp.ComponentIndex[cID] = append(cp.ComponentIndex[cID], arch)
	cp.ArchetypeIndex[arch] = append(cp.ArchetypeIndex[arch], e)
	cp.ArchetypeDefinitions[arch] = append(cp.ArchetypeDefinitions[arch], cID)
	// add component to the world
	AddToWorld(comp)

}
type ComponentSlice[T any] struct {
	slice map[ArchetypeID][]T
}
