package ecsgo

type World struct {
	creatures *ComponentSlice[Creature]
}

type ComponentSlice[T any] struct {
	slice []T
}
