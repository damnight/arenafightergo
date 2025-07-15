package ecsgo

import (
	"fmt"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Coordinator struct {

	// Managers
	em *EntityManager
	cm *ComponentManager
	am *ArchetypeManager
	sm *SystemsManager

	renderList  []EntityID
	spriteSheet *SpriteSheet
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
	}

	return co, nil

}

func (co *Coordinator) renderLevel(screen *ebiten.Image, g *Game) {
	start := time.Now()
	co.UpdateRenderList()
	update_render := time.Since(start)

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

	start = time.Now()
	for _, e := range co.renderList {
		compIDs := co.em.EntityIndex[e]
		pos := co.cm.GetComponentByID(e, compIDs, PositionType)
		sp := co.cm.GetComponentByID(e, compIDs, SpriteType)

		if position, ok := pos.(Position); ok {
			x, y := position.x, position.y

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
			if sprite, ok2 := sp.(Sprite); ok2 {
				target.DrawImage(sprite.img[0], op)
			}
		}
	}
	renderloop := time.Since(start)

	start = time.Now()
	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.CamScale), float64(g.CamScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}
	scaleLaterdraw := time.Since(start)

	fmt.Printf("| Update Randerlist: %v | Render Loop: %v | Scale Later draw: %v |\n", update_render, renderloop, scaleLaterdraw)

}

func (co *Coordinator) UpdateRenderList() {
	arch := co.FindRenderArchetype()
	renderList := co.am.ArchetypeIndex[arch]
	co.renderList = renderList
}

func (co *Coordinator) FindRenderArchetype() ArchetypeID {
	sig := []ComponentTypeID{}
	sig = append(sig, PositionType)
	sig = append(sig, SpriteType)
	slices.Sort(sig)

	var renderArch ArchetypeID
	for archID, signature := range co.am.ArchetypeDefinitions {
		if slices.Equal(signature, sig) {
			renderArch = archID
			break
		}
	}

	return renderArch
}

func (co *Coordinator) AddEntity(comps []IComponent) (EntityID, error) {
	// create entity
	e, err := co.em.CreateEntity()
	if err != nil {
		return 0, err
	}

	// Register and add components to maps and slices respectively
	compIDList := co.cm.RegisterComponents(e, comps)

	// EntityIndex[EntityID][]ComponentID
	co.em.EntityIndex[e] = compIDList

	// check if archetype exists for this set of Components, add or create accordingly
	slices.Sort(compIDList)
	arch := co.am.GetSetArchetype(comps, co.cm)

	//map[ComponentID][]ArchetypeID
	for _, c := range comps {
		id := co.cm.ComponentDefinitions[&c]
		co.cm.ComponentIndex[id] = append(co.cm.ComponentIndex[id], arch)
	}

	//map[ArchetypeID][]EntityID
	co.am.ArchetypeIndex[arch] = append(co.am.ArchetypeIndex[arch], e)

	return e, nil

}

func (co *Coordinator) CreateTile(x, y, z float64, tileType SpriteID) EntityID {

	compList := []IComponent{}
	compList = append(compList, Position{x: x, y: y, z: z})
	compList = append(compList, Sprite{spriteID: tileType, img: co.spriteSheet.slice[tileType]})

	tile, err := co.AddEntity(compList)
	if err != nil {
		return 0
	}

	return tile
}

//
//ype World struct {
//       index map[ArchetypeID]*ComponentSlice[ArchetypeID]
//
//
//unc CreateWorld() *World {
//       return &World{index: make(map[ArchetypeID]*ComponentSlice[ArchetypeID])}
//
//
//unc (w *World) addToWorld(arch ArchetypeID, comps []IComponent) {
//       cs := CreateComponentSlice()
//       w.index[arch] = cs
//       cs.addComponents(arch, comps)
//
//
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
