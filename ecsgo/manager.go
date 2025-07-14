package ecsgo

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type Coordinator struct {

	// Managers
	em *EntityManager
	cm *ComponentManager
	am *ArchetypeManager
	sm *SystemsManager

	//world       *World
	renderList  *[]IComponent
	spriteSheet *SpriteSheet
}

func (co *Coordinator) GetTileSprites(tileType SpriteID) *[]*Sprite {

	tileSprites := co.spriteSheet.slice[tileType]
	return &tileSprites
}

func NewCoordinator() (*Coordinator, error) {
	sheet, err0 := LoadSpriteSheet(64)
	em, err1 := NewEntityManager()
	cm, err2 := NewComponentManager()
	am, err3 := NewArchetypeManager()
	sm, err4 := NewSystemsManager()

	if err0 != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return &Coordinator{}, fmt.Errorf("Manager initialization error: %s\n%s\n%s\n%s\n%s",
			err0, err1, err2, err3, err4)
	}

	co := &Coordinator{
		em:          em,
		cm:          cm,
		am:          am,
		sm:          sm,
		spriteSheet: sheet,
		world:       CreateWorld(),
		renderList:  CreateComponentSlice(),
	}

	return co, nil

}

func (co *Coordinator) renderLevel(screen *ebiten.Image, g *Game) {
	co.UpdateRenderList()
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

	for i, _ := range co.renderList.data {

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

func (co *Coordinator) UpdateRenderList() {
	arch := co.FindRenderArchetype()
	renderList := co.getComponentSlice(arch)
	co.renderList = renderList
}
func (co *Coordinator) getComponentSlice(arch ArchetypeID) *ComponentSlice[ArchetypeID] {
	return co.world.index[arch]
}

func (co *Coordinator) FindRenderArchetype() ArchetypeID {
	comps := []IComponent{}
	comps = append(comps, Position{})
	comps = append(comps, Sprite{})

	arch := co.am.GetSetArchetype(comps, co.cm)
	return arch
}

func (co *Coordinator) AddEntity(comps []IComponent) (*EntityID, error) {
	// create entity
	e, err := co.em.CreateEntity()
	if err != nil {
		return nil, err
	}

	// create componentIDs
	compList := []ComponentID{}
	for _, c := range comps {
		id := co.cm.CreateComponentID()
		// ComponentDefinitions[*IComponent]ComponentID
		co.cm.ComponentDefinitions[c.Type()] = id
		compList = append(compList, id)
	}

	// EntityIndex[EntityID][]ComponentID
	co.em.EntityIndex[e] = compList

	// check if archetype exists for this set of Components, add or create accordingly
	slices.Sort(compList)
	arch := co.am.GetSetArchetype(comps)

	//map[ComponentID][]ArchetypeID
	for _, c := range comps {
		cp.ComponentIndex[c.Type()] = append(cp.ComponentIndex[c.Type()], arch)
	}

	//map[ArchetypeID][]EntityID
	co.ArchetypeIndex[arch] = append(co.ArchetypeIndex[arch], e)

	// add entity under archetype to data in a componentslice, with index
	//cp.world.addToWorld(arch, comps)
	co.world.addToWorld(arch, comps)
	return &e

}

func (co *Coordinator) CreateTile(x, y, z float64, tileType SpriteID) *EntityID {

	compList := []IComponent{}
	compList = append(compList, Position{x: x, y: y, z: z})
	compList = append(compList, Sprite{spriteID: spriteID, img: cp.spriteSheet.slice[spriteID]})

	tile := co.AddEntity(compList)

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
	data map[ArchetypeID][]IComponent
}

func (cs *ComponentSlice[T]) addComponents(arch ArchetypeID, comps []IComponent) {
	cs.data[arch] = append(cs.data[arch], comps...)

}

func CreateComponentSlice() *ComponentSlice[ArchetypeID] {
	return &ComponentSlice[ArchetypeID]{
		data: make(map[ArchetypeID][]IComponent),
	}

}

//ype ComponentSlice[T any] struct {
//       data  map[int][]IComponent
//       index map[ArchetypeID]int
//
//
//unc (cs *ComponentSlice[T]) addComponents(arch ArchetypeID, comps []IComponent) {
//       cs.data[cs.index[arch]] = append(cs.data[cs.index[arch]], comps)
//
//
//unc CreateComponentSlice() *ComponentSlice[ArchetypeID] {
//       return &ComponentSlice[ArchetypeID]{
//       	data:  make(map[int][]IComponent),
//       	index: make(map[ArchetypeID]int),
//       }
//
//
