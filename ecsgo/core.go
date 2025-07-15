package ecsgo

import (
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"time"
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
	if em.entityIDCounter < MAX_ENTITIES {
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
			cType := cm.ComponentIDType[cID]

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
	ComponentIDIndex      map[ComponentID]*IComponent
	ComponentsByTypeIndex map[ComponentTypeID][]ComponentID
	ComponentDefinitions  map[*IComponent]ComponentID
	ComponentIDType       map[ComponentID]ComponentTypeID
	cf                    *ComponentField
}

func NewComponentManager() (*ComponentManager, error) {
	cm := &ComponentManager{
		ComponentIndex:        make(map[ComponentID][]ArchetypeID),
		ComponentsByTypeIndex: make(map[ComponentTypeID][]ComponentID),
		ComponentDefinitions:  make(map[*IComponent]ComponentID),
		ComponentIDIndex:      make(map[ComponentID]*IComponent),
		ComponentIDType:       make(map[ComponentID]ComponentTypeID),
	}

	NewComponentField(cm)
	return cm, nil
}

func (cm *ComponentManager) CreateComponentID() ComponentID {
	id := atomic.AddUint64(&cm.componentIDCounter, 1)
	return ComponentID(id)
}

func (cm *ComponentManager) RegisterComponents(e EntityID, comps []IComponent) []ComponentID {
	var compIDList []ComponentID

	for _, c := range comps {
		id := cm.CreateComponentID()
		// Register component with ID
		cType := c.Type()
		cm.ComponentIDIndex[id] = &c
		cm.ComponentDefinitions[&c] = id
		cm.ComponentsByTypeIndex[cType] = append(cm.ComponentsByTypeIndex[cType], id)
		cm.ComponentIDType[id] = cType
		compIDList = append(compIDList, id)
	}

	for _, c := range comps {
		cm.cf.AddComponent(e, c)
	}

	return compIDList

}

func (cm *ComponentManager) GetComponentByID(e EntityID, comps []ComponentID, cType ComponentTypeID) IComponent {
	for _, cID := range comps {
		start := time.Now()
		cTypeX := cm.ComponentIDType[cID]
		finish := time.Since(start)
		fmt.Printf("|>>| TypeFromComponentID: %v |\n", finish)
		if cTypeX == cType {
			c, err1 := cm.cf.GetComponent(e, cID, cTypeX)
			if err1 != nil {
				return nil
			}
			return c
		}
	}
	return nil
}

type ComponentField struct {
	positions *ComponentSlice[Position]
	sprites   *ComponentSlice[Sprite]
}

func NewComponentField(cm *ComponentManager) {
	cm.cf = &ComponentField{}

	NewComponentSlice(PositionType, cm)
	NewComponentSlice(SpriteType, cm)
}

func (cf *ComponentField) AddComponent(e EntityID, c IComponent) {
	switch c.Type() {
	case PositionType:
		cf.positions.addComponent(e, c.(Position))
	case SpriteType:
		cf.sprites.addComponent(e, c.(Sprite))
	default:
		return
	}
}

func (cf *ComponentField) RemoveComponent(e EntityID, c IComponent) {
	switch c.Type() {
	case PositionType:
		cf.positions.removeComponent(e)
	case SpriteType:
		cf.sprites.removeComponent(e)
	default:
		return
	}
}

func (cf *ComponentField) GetComponent(e EntityID, cID ComponentID, cType ComponentTypeID) (IComponent, error) {
	var err error
	var c IComponent

	switch cType {
	case PositionType:
		c, err = cf.positions.getComponent(e)
	case SpriteType:
		c, err = cf.sprites.getComponent(e)
	default:
		return c, fmt.Errorf("Default type error, no components found")
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
	case PositionType:
		cm.cf.positions = &ComponentSlice[Position]{data: []Position{}, entityMap: make(map[EntityID]uint)}
		return nil
	case SpriteType:
		cm.cf.sprites = &ComponentSlice[Sprite]{data: []Sprite{}, entityMap: make(map[EntityID]uint)}
		return nil
	default:
		return fmt.Errorf("Can´t Create NilType ComponentSlice")

	}
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
