package ecsgo

import (
	"fmt"
)

type CreationFactory struct {
	// EntityID -> []ComponentID
	// ComponentID -> []ArchetypID
	// ArchetypeID -> []EntityID
	EntityIndex    map[EntityID][]ComponentID
	ComponentIndex map[ComponentID][]ArchetypeID
	ArchetypeIndex map[ArchetypeID][]EntityID

	ArchetypeDefinitions map[ArchetypeID][]ComponentID

	spriteSheet *SpriteSheet
}

func NewCreationFactory() (CreationFactory, error) {
	cf := CreationFactory{}
	sheet, err := LoadSpriteSheet(64)
	if err != nil {
		return CreationFactory{}, fmt.Errorf("failed to load embedded spritesheet: %s", err)
	}

	cf.spriteSheet = sheet

	return cf, nil

}

func (cf *CreationFactory) CreatePlayer1() {

	//	//create player character archetype
	//	archPlayerChar := CreateArchetypeID()
	//
	//	//create player character composition
	//	compList := []IComponent{}
	//
	//	compList = append(compList, Position{0.0, 0.0, 0.0})
	//	compList = append(compList, Velocity{0.0, 0.0})
	//	compList = append(compList, Health{100})
	//	compList = append(compList, BaseSpeed{30})
	//	compList = append(compList, cf.spriteSheet.creature[Knight])
	//
	//	//add archetype
	//	cf.ArchetypeDefinitions[archPlayerChar] = append(cf.ArchetypeDefinitions[archPlayerChar], compList)
	//
	//	for _, comp := range compList {
	//		cf.ComponentIndex[comp] = archPlayerChar
	//	}
	//
	//	//create player character entity
	//	e := CreateEnCreateEntity()
	//
	//	//add entity
	//	cf.EntityIndex[e] = compList

}
