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
	renderList  *ComponentSlice[ArchetypeID]
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
		renderList: CreateComponentSlice(),
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
func (cp *ComponentManager) renderLevel(screen *ebiten.Image, g *Game) {
	cp.UpdateRenderList()
	op := &ebiten.DrawImageOptions{}
	padding := float64(g.CurrentLevel.TileSize) * g.CamScale
	cx, cy := float64(g.Width/2), float64(g.Height/2)

	scaleLater := g.CamScale > 1
	target := screen
	scale := g.CamScale

	// When zooming in, tiles can have slight bleeding edges.
	// To avoid them, render the result on an Offscreen first and then scale it later.
	if scaleLater {
		if g.Offscreen != nil {
			if g.Offscreen.Bounds().Size() != screen.Bounds().Size() {
				g.Offscreen.Deallocate()
				g.Offscreen = nil
			}
		}
		if g.Offscreen == nil {
			s := screen.Bounds().Size()
			g.Offscreen = ebiten.NewImage(s.X, s.Y)
		}
		target = g.Offscreen
		target.Clear()
		scale = 1
	}

	for i, _ := range cp.renderList.data {

		x := 1
		y := i
		xi, yi := g.cartesianToIso(float64(x), float64(y))

		// Skip drawing tiles that are out of the screen.
		drawX, drawY := ((xi-g.CamX)*g.CamScale)+cx, ((yi+g.CamY)*g.CamScale)+cy
		if drawX+padding < 0 || drawY+padding < 0 || drawX > float64(g.Width) || drawY > float64(g.Height) {
			continue
		}

		op.GeoM.Reset()
		// Move to current isometric position.
		op.GeoM.Translate(xi, yi)
		// Translate camera position.
		op.GeoM.Translate(-g.CamX, g.CamY)
		// Zoom.
		op.GeoM.Scale(scale, scale)
		// Center.
		op.GeoM.Translate(cx, cy)

		screen.DrawImage(&ebiten.Image{}, op)
	}
	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.CamScale), float64(g.CamScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}
}

func (cp *ComponentManager) UpdateRenderList() {
	arch := cp.FindRenderArchetype()
	renderList := cp.getComponentSlice(arch)
	cp.renderList = renderList
}

func (cp *ComponentManager) getComponentSlice(arch ArchetypeID) *ComponentSlice[ArchetypeID] {
	return cp.world.index[arch]

}

func (cp *ComponentManager) FindRenderArchetype() ArchetypeID {
	comps := []IComponent{}
	comps = append(comps, &Position{})
	comps = append(comps, Default)

	arch := cp.CheckForArchetype(comps)
	return arch
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

	//map[ComponentID][]ArchetypeID
	for _, c := range compList {
		cp.ComponentIndex[c] = append(cp.ComponentIndex[c], arch)
	}

	//map[ArchetypeID][]EntityID
	cp.ArchetypeIndex[arch] = append(cp.ArchetypeIndex[arch], e)

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
