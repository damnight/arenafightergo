package ecsgo

import "sync/atomic"

type EntityID uint64
type ArchetypeID uint64
type ComponentID uint64

var entityIDCounter uint64
var archetypeIDCounter uint64
var componentIDCounter uint64

func CreateEntity() EntityID {
	id := atomic.AddUint64(&entityIDCounter, 1)
	return EntityID(id)
}

func CreateArchetypeID() ArchetypeID {
	id := atomic.AddUint64(&archetypeIDCounter, 1)
	return ArchetypeID(id)
}

func CreateComponentID() ComponentID {
	id := atomic.AddUint64(&componentIDCounter, 1)
	return ComponentID(id)
}
