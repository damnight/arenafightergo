package ecsgo

import (
	"errors"
	"fmt"
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

func (em *EntityManager) DestroyEntity(e EntityID, cm *ComponentManager) error {
	if e < MAX_ENTITIES {
		compIDs := em.EntityIndex[e]
		for _, cID := range compIDs {
			cType, err := cm.TypeFromComponentID(cID)
			if err != nil {
				return err
			}

			c, err := cm.cf.GetComponent(e, cID, cType)
			if err != nil {
				return err
			}
			cm.cf.RemoveComponent(e, c)
		}
		// TODO: deregister from archetypes

		delete(em.EntityIndex, e)

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
	return nil
}
func (em *EntityManager) GetEntity(cID ComponentID) (EntityID, error) {
	for e, compID := range em.EntityIndex {
		if slices.Contains(compID, cID) {
			return e, nil
		}
	}
	return 0, fmt.Errorf("No Entity found for Component: %v", cID)
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
	cf                    *ComponentField
}

func NewComponentManager() (*ComponentManager, error) {
	cm := &ComponentManager{
		ComponentIndex:        make(map[ComponentID][]ArchetypeID),
		ComponentsByTypeIndex: make(map[ComponentTypeID][]ComponentID),
		ComponentDefinitions:  make(map[*IComponent]ComponentID),
	}

	cm.cf = NewComponentField(cm)
	return cm, nil
}

func (cm *ComponentManager) CreateComponentID() ComponentID {
	id := atomic.AddUint64(&cm.componentIDCounter, 1)
	return ComponentID(id)
}

func (cm *ComponentManager) RegisterComponents(e EntityID, comps []IComponent) {
	for _, c := range comps {
		cm.cf.AddComponent(e, c)
	}
}

func (cm *ComponentManager) TypeFromComponentID(cID ComponentID) (ComponentTypeID, error) {
	for cType, compIDs := range cm.ComponentsByTypeIndex {
		if slices.Contains(compIDs, cID) {
			return cType, nil
		}
	}
	return VoidType, fmt.Errorf("ComponentID had no Type Entry")
}

type ComponentField struct {
	positions *ComponentSlice[Position]
}

func NewComponentField(cm *ComponentManager) *ComponentField {
	cf := &ComponentField{}
	for i := uint(1); i < uint(MAX_COMPONENTTYPE_ID)-1; i++ {
		NewComponentSlice(ComponentTypeID(i), cm)
	}

	return cf

}

func (cf *ComponentField) AddComponent(e EntityID, c IComponent) {
	switch c.Type() {
	case VoidType:
		return
	case PositionType:
		cf.positions.addComponent(e, c.(Position))
	}
}

func (cf *ComponentField) RemoveComponent(e EntityID, c IComponent) {
	switch c.Type() {
	case VoidType:
		return
	case PositionType:
		cf.positions.removeComponent(e)
	}
}

func (cf *ComponentField) GetComponent(e EntityID, cID ComponentID, cType ComponentTypeID) (IComponent, error) {
	var err error
	var c IComponent

	switch cType {
	case PositionType:
		c, err = cf.positions.getComponent(e)
	}

	if err != nil {
		return c, fmt.Errorf("Component Not Found")
	}
	return c, err
}

type ComponentSlice[T any] struct {
	data      []T               // arrayslice of all components
	entityMap map[EntityID]uint // {Entity, data.index of component}}
}

func NewComponentSlice(cType ComponentTypeID, cm *ComponentManager) error {
	switch cType {
	case VoidType:
		return fmt.Errorf("Can´t Create NilType ComponentSlice")
	case PositionType:
		cm.cf.positions = &ComponentSlice[Position]{data: []Position{}, entityMap: make(map[EntityID]uint)}
		return nil
	}

	return nil

}

func (cs *ComponentSlice[T]) addComponent(e EntityID, component T) {
	if _, exists := cs.entityMap[e]; exists {
		cs.data[cs.entityMap[e]] = component
	} else {
		// add new
		cs.data = append(cs.data, component)
		cs.entityMap[e] = uint(cs.Length()) - 1
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

func (cs *ComponentSlice[T]) removeComponent(e EntityID) {
	if idx, exists := cs.entityMap[e]; exists {
		// remove existing component, swap and delete
		lastIdx := uint(len(cs.data)) - 1
		cs.data[idx] = cs.data[lastIdx]
		cs.data = cs.data[:lastIdx]

		// update entityMap
		// now the old entityMap which pointed at data[lastIdx] is out of bounds, and the entityMap for e points at data[idx]

		// reverse lookup, find the entity in the entityMap which still points to lastIdx
		var swappedEntity EntityID
		for entity, dataIdx := range cs.entityMap {
			if dataIdx == lastIdx {
				swappedEntity = entity
				break
			}
		}

		// give the swapped entity the correct data index, which is the old deleted one
		cs.entityMap[swappedEntity] = idx
		// delete the old entityMap entry, e still points at idx otherwise
		delete(cs.entityMap, e)

	} else {
		// TODO: error handling
	}

}

func (cs *ComponentSlice[T]) Length() int {
	return len(cs.data)
}
