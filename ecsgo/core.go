package ecsgo

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync/atomic"
)

type EntityID uint64
type ArchetypeID uint64
type ComponentID uint64

const MAX_ENTITIES = 100000

type EntityManager struct {
	// EntityID -> []ComponentID
	EntityIndex     map[EntityID][]ComponentID
	entityIDCounter uint64
}

func NewEntityManager() (*EntityManager, error) {
	// Initialize the queue with all possible entity IDs
	em := &EntityManager{
		EntityIndex: make(map[EntityID][]ComponentID),
	}

	return em, nil
}

func (em *EntityManager) CreateEntity() (EntityID, error) {
	if !(em.entityIDCounter < MAX_ENTITIES) {
		id := atomic.AddUint64(&em.entityIDCounter, 1)
		return EntityID(id), nil
	} else {
		return MAX_ENTITIES + 1, fmt.Errorf("Too many Entities!")
	}
}

func (em *EntityManager) DestroyEntity(entity EntityID) error {
	if entity < MAX_ENTITIES {
		delete(em.EntityIndex, entity)
	} else {
		return fmt.Errorf("Entity out of range!")
	}
	// assert(entity < MAX_ENTITIES && "Entity out of range.");
	//
	// // Invalidate the destroyed entity's signature
	// mSignatures[entity].reset();
	//
	// // Put the destroyed ID at the back of the queue
	// Registry.push(entity);
	// --mLivingEntityCount;
}

// Archetype Manager
type ArchetypeManager struct {
	// ArchetypeID -> []EntityID
	ArchetypeIndex     map[ArchetypeID][]EntityID
	archetypeIDCounter uint64
	// Signature
	ArchetypeDefinitions map[ArchetypeID][]ComponentTypeID
}

func NewArchetypeManager() (*ArchetypeManager, error) {
	am := &ArchetypeManager{
		ArchetypeIndex:       make(map[ArchetypeID][]EntityID),
		ArchetypeDefinitions: make(map[ArchetypeID][]ComponentTypeID),
	}

	return am, nil
}

func (am *ArchetypeManager) CreateArchetypeID() ArchetypeID {
	id := atomic.AddUint64(&am.archetypeIDCounter, 1)
	return ArchetypeID(id)
}

func (am *ArchetypeManager) GetSetArchetype(comps []IComponent, cm *ComponentManager) ArchetypeID {
	// get all assiciated ComponentIDs by ArchetypeDefinition
	// get all associated Archetypes by ComponentIndex
	archList := []ArchetypeID{}
	compTypeList := []ComponentTypeID{}

	// retrieve all matching archetypes
	for _, c := range comps {
		id := cm.ComponentDefinitions[&c]
		compTypeList = append(compTypeList, c.Type())

		archetypeIDs := cm.ComponentIndex[id]
		for _, archID := range archetypeIDs {
			// for loop over archList, if not already found, add to archList
			if slices.Index(archList, archID) < 0 {
				archList = append(archList, archID)
			}
		}
	}

	equality := false
	slices.Sort(compTypeList)
	// check if archetype in ArchetypeDefinitions
	for _, arch := range archList {
		archtypeComponents := am.ArchetypeDefinitions[arch]
		// slices.Equal compares idx by idx, so slice needs to be sorted, fortunately ArchID is a uint
		slices.Sort(archtypeComponents)

		equality = slices.Equal(archtypeComponents, compTypeList)
		if equality {
			return arch
		}
	}
	// only executes if no match was found, thus creating a new archetype
	archID := am.CreateArchetypeID()
	am.ArchetypeDefinitions[archID] = compTypeList

	return archID
}

// do I actually need a delete function?
//
//	func (am *ArchetypeManager) DeleteUnsusedArchetypes() {
//	      // run this periodically as garbage collectio
//
// Component Manager
type ComponentManager struct {
	componentIDCounter uint64
	// ComponentID -> []ArchetypID
	ComponentIndex        map[ComponentID][]ArchetypeID
	ComponentsByTypeIndex map[ComponentTypeID][]ComponentID
	ComponentDefinitions  map[*IComponent]ComponentID
	cs                    []*ComponentSlice
}

func NewComponentManager() (*ComponentManager, error) {
	cm := &ComponentManager{
		ComponentIndex:        make(map[ComponentID][]ArchetypeID),
		ComponentsByTypeIndex: make(map[ComponentTypeID][]ComponentID),
		ComponentDefinitions:  make(map[*IComponent]ComponentID),
	}
	return cm, nil
}

func (cm *ComponentManager) CreateComponentID() ComponentID {
	id := atomic.AddUint64(&cm.componentIDCounter, 1)
	return ComponentID(id)
}

func (cm *ComponentManager) RegisterComponents(e EntityID, comps []IComponent) {
	for _, c := range comps {
		cType := c.Type()
		if _, exists := cm.cs[cType]; exists {
			cm.cs[cType].addComponent(e, c)

		} else {
			// get c type
			cm.cs[cType] = NewComponentSlice(cType)
		}

	}

}

type ComponentSlice[T any] struct {
	data      []T               // arrayslice of all components
	entityMap map[EntityID]uint // {Entity, data.index of component}}
}

func NewComponentSlice[T any](cType ComponentTypeID) (ComponentSlice[T], error) {
	switch cType {
	case NilType:
		return nil, fmt.Errorf("Can´t Create NilType ComponentSlice")
	case PositionType:
		return ComponentSlice[Position]{data: []Position{}, entityMap: make(map[EntityID]uint)}, nil
	}

	return nil, fmt.Errorf("Switch Statement didn´t find Type")

}

func (cs *ComponentSlice[T]) addComponent(e EntityID, component T) {
	if _, exists := cs.entityMap[e]; exists {

		cs.data[cs.entityMap[e]] = component
	} else {
		// add new
		cs.data = append(cs.data, component)
		cs.entityMap[e] = len(cs.data) - 1
	}
}

func (cs *ComponentSlice[T]) getComponent(e EntityID) (T, error) {
	if idx, exists := cs.entityMap[e]; exists {
		// return component
		return cs.data[idx], nil

	} else {
		// return error
		var undefined T
		return undefined, errors.New("Component not found")
	}
}

func (cs *ComponentSlice[T]) Length() int {
	return len(cs.data)
}
