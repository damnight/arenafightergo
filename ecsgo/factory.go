package ecsgo

type CreationFactory struct {
}

func (cp *ComponentManager) CreateTile(x, y, z float64, tileType SpriteID) *Tile {
	
	// create entity
	e := CreateEntity()

	tile := Tile{
	// create object components
	position : &Position{},
	sprites : GetTileSprites(tileType)
	}

	cp.AddEntity(e, tile)



	// add obj as component

	cp.AddComponent(tile)

	return tile
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
