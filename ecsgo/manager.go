package ecsgo

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type ComponentManager struct {
	// EntityID -> []ComponentID
	// ComponentID -> []ArchetypID
	// ArchetypeID -> []EntityID
	EntityIndex          map[EntityID][]ComponentID
	ComponentIndex       map[ComponentID][]ArchetypeID
	ArchetypeIndex       map[ArchetypeID][]EntityID
	ArchetypeDefinitions map[ArchetypeID][]ComponentID
	ComponentDefinitions map[*IComponent]ComponentID

	world       *World
	renderList  *Renderable
	spriteSheet *SpriteSheet
}

func (cp *ComponentManager) GetTileSprites(tileType SpriteID) *[]*Sprite {

	tileSprites := cp.spriteSheet.slice[tileType]
	return &tileSprites
}

func NewComponentManager() (*ComponentManager, error) {
	cp := &ComponentManager{
		EntityIndex:          make(map[EntityID][]ComponentID),
		ComponentIndex:       make(map[ComponentID][]ArchetypeID),
		ArchetypeIndex:       make(map[ArchetypeID][]EntityID),
		ArchetypeDefinitions: make(map[ArchetypeID][]ComponentID),
		ComponentDefinitions: make(map[*IComponent]ComponentID),

		world:      CreateWorld(),
		renderList: &Renderable{},
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
func (cp *ComponentManager) DrawTileSprites(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	for _, r := range cp.spriteSheet.slice {
		for _, s := range r {
			screen.DrawImage(s.img, options)
		}
	}
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

func (cp *ComponentManager) AddEntity(comps []IComponent) *EntityID {
	// create entity
	e := CreateEntity()
	// create componentIDs
	compList := []ComponentID{}
	for _, c := range comps {
		id := CreateComponentID()
		// ComponentDefinitions[*IComponent]ComponentID
		cp.ComponentDefinitions[&c] = id
		compList = append(compList, id)
	}

	// EntityIndex[EntityID][]ComponentID
	cp.EntityIndex[e] = compList

	// check if archetype exists for this set of Components, add or create accordingly
	arch := cp.CheckForArchetype(comps)

	// add entity under archetype to data in a componentslice, with index
	cp.world.addToWorld(arch, comps)

	return &e

}

func (cp *ComponentManager) CreateTile(x, y, z float64, tileType SpriteID) *EntityID {

	compList := []IComponent{}
	compList = append(compList, Position{x: x, y: y, z: z})
	compList = append(compList, tileType)

	tile := cp.AddEntity(compList)

	return tile
}

type World struct {
	index map[ArchetypeID]*ComponentSlice[ArchetypeID]
}

func CreateWorld() *World {
	return &World{index: make(map[ArchetypeID]*ComponentSlice[ArchetypeID])}
}

func (w *World) addToWorld(arch ArchetypeID, comps []IComponent) {
	cs := CreateComponentSlice()
	w.index[arch] = cs
	cs.addComponents(arch, comps)
}

type ComponentSlice[T any] struct {
	data  map[int][]IComponent
	index map[ArchetypeID]int
}

func (cs *ComponentSlice[T]) addComponents(arch ArchetypeID, comps []IComponent) {
	cs.data[cs.index[arch]] = append(cs.data[cs.index[arch]], comps)
}

func CreateComponentSlice() *ComponentSlice[ArchetypeID] {
	return &ComponentSlice[ArchetypeID]{
		data:  make(map[int][]IComponent),
		index: make(map[ArchetypeID]int),
	}

}
