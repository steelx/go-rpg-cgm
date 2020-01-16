package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	log "github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/tilepix"
	"image/color"
)

type GameMap struct {
	x, y float64

	// To track the camera position
	CamX, CamY float64

	MapInfo MapInfo
	sprites map[string]*pixel.Sprite

	mTileSprite   pixel.Sprite
	Width, Height float64

	Tiles        []*pixel.Batch
	tilesIndices map[string]int
	tilesCounter int

	hideDecorationTile    []bool //(tileX)+(tileY * 100)
	bypassBlockedTile     map[[2]float64]bool
	TileWidth, TileHeight float64
	blockingTileGID       tilepix.GID
	Canvas                *pixelgl.Canvas
	renderLayer           int

	Actions        map[string]func(gMap *GameMap, entity *Entity, x, y float64)
	TriggerTypes   map[string]Trigger
	Triggers       map[[2]float64]Trigger
	OnWakeTriggers map[string]Trigger

	Entities     []*Entity
	NPCs         []*Character
	NPCbyId      map[string]*Character
	MarkReRender bool
}

func MapCreate(mapInfo MapInfo) *GameMap {
	m := &GameMap{
		MapInfo:           mapInfo,
		bypassBlockedTile: make(map[[2]float64]bool),
	}

	m.NPCbyId = make(map[string]*Character, 0)
	m.Entities = make([]*Entity, 0)

	m.Height = float64(mapInfo.Tilemap.Height)
	m.Width = float64(mapInfo.Tilemap.Width)

	m.TileWidth = float64(mapInfo.Tilemap.TileWidth)
	m.TileHeight = float64(mapInfo.Tilemap.TileHeight)

	//Bottom left corner of the map, since pixel starts at 0, 0
	m.x = m.TileWidth
	m.y = m.TileHeight

	m.Canvas = pixelgl.NewCanvas(m.MapInfo.Tilemap.Bounds())
	m.setTiles()
	m.setBlockingTileInfo()
	m.createTriggersFromMapInfo()
	return m
}

func (m *GameMap) createTriggersFromMapInfo() {
	m.Actions = make(map[string]func(gMap *GameMap, entity *Entity, x, y float64))
	for name, def := range m.MapInfo.Actions {
		//def.Id = RunScript
		action := RunScript(def.Script)
		m.Actions[name] = action
	}

	//Create the Trigger types from the Action def
	m.TriggerTypes = make(map[string]Trigger)
	for k, v := range m.MapInfo.TriggerTypes {
		m.TriggerTypes[k] = Trigger{
			OnEnter: m.Actions[v.OnEnter],
			OnExit:  m.Actions[v.OnExit],
			OnUse:   m.Actions[v.OnUse],
		}
	}

	m.Triggers = make(map[[2]float64]Trigger)
	for _, v := range m.MapInfo.Triggers {
		//we take Tile XY and set as map x, y cords
		x, y := m.GetTileIndex(v.X, v.Y)
		m.Triggers[[2]float64{x, y}] = m.TriggerTypes[v.Id]
	}

	m.OnWakeTriggers = make(map[string]Trigger)
	for key, v := range m.MapInfo.OnWake {
		addNPC := LIST[key](m, v.X, v.Y)
		addNPC(Characters[v.Id](m))
	}

	m.hideDecorationTile = make([]bool, m.MapInfo.Tilemap.Width*m.MapInfo.Tilemap.Height)

}

//SetHiddenTileVisible will set hideDecorationTile to true
// but actually we are setting decoration tile to hide
//so we can show Background layer tile below
func (m *GameMap) SetHiddenTileVisible(tileX, tileY int) {
	//assuming Height * Width are same e.g. 100x100
	m.MarkReRender = true
	m.hideDecorationTile[tileX+(tileY*m.MapInfo.Tilemap.Height)] = true
}

func (m *GameMap) setBlockingTileInfo() {
	for _, tile := range m.MapInfo.Tilemap.Tilesets {
		if tile.Name == "collision_px" {
			m.blockingTileGID = tile.FirstGID
			break
		}
	}
}
func (m *GameMap) ClearAllEntities() {
	m.Entities = make([]*Entity, 0)
}

func (m GameMap) GetEntityAtPos(x, y float64) *Entity {
	for _, e := range m.Entities {
		if e.TileX == x && e.TileY == y {
			return e
		}
	}
	return nil
}

//IsBlockingTile check's X, Y cords on collision map layer
// if ID is not 0, tile exists on X, Y we return true
func (m GameMap) IsBlockingTile(tileX, tileY int) bool {
	if (tileX + tileY*int(m.Width)) <= 0 {
		return true //we dont let him go out of map
	}

	if x, y := m.GetTileIndex(float64(tileX), float64(tileY)); m.bypassBlockedTile[[2]float64{x, y}] {
		return false
	}
	tile := m.MapInfo.Tilemap.TileLayers[m.MapInfo.CollisionLayer].DecodedTiles[tileX+(tileY*int(m.Width))]
	return !tile.IsNil() || tile.ID != 0
}

func (m *GameMap) setTiles() {
	batches := make([]*pixel.Batch, 0)
	batchIndices := make(map[string]int)
	batchCounter := 0

	// Load the sprites
	sprites := make(map[string]*pixel.Sprite)
	for _, tileset := range m.MapInfo.Tilemap.Tilesets {
		if _, alreadyLoaded := sprites[tileset.Image.Source]; !alreadyLoaded {
			sprite, pictureData := utilz.LoadSprite(tileset.Image.Source)
			sprites[tileset.Image.Source] = sprite
			batches = append(batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
			batchIndices[tileset.Image.Source] = batchCounter
			batchCounter++
		}
	}
	m.Tiles = batches
	m.tilesIndices = batchIndices
	m.tilesCounter = batchCounter
	m.sprites = sprites
}

//Cam to Tile : GoToTile pan camera to given coordinates
func (m *GameMap) GoToTile(x, y float64) {
	tileX, tileY := m.GetTileIndex(x, y)
	x = tileX - m.TileWidth/2
	y = tileY - m.TileHeight/2
	m.Goto(x, y)
}

func (m *GameMap) Goto(x, y float64) {
	m.CamX = x
	m.CamY = y
}

//GetTileIndex will take TileX, TileY and return exact MAP cords
//e.g. 35, 22 will return cords on map x 400, y 1300
func (m GameMap) GetTileIndex(tileX, tileY float64) (x, y float64) {
	tileY = m.Height - tileY //make count Y from top (Tiled app starts from top)
	x = m.x + (tileX * m.TileWidth)
	y = m.y + (tileY * m.TileHeight)
	return
}

func (m GameMap) GetTilePositionAtFeet(x, y, charW, charH float64) pixel.Vec {
	tileX, tileY := m.GetTileIndex(x, y)
	x = tileX - (charW / 2)
	y = tileY - (charH / 2) - 5
	return pixel.V(x, y)
}

//func (m GameMap) DrawAll(target pixel.Target, clearColour color.Color, mat pixel.Matrix) {
//	m.MapInfo.Tilemap.DrawAll(target, clearColour, mat)
//}

//DrawAfter will render the callback function after given layer index
// uses pixelgl Canvas instead of Win to render
func (m GameMap) DrawAfter(callback func(canvas *pixelgl.Canvas, layer int)) error {
	// Draw tiles
	target, mat := globals.Global.Win, pixel.IM

	if m.Canvas == nil {
		m.Canvas = pixelgl.NewCanvas(m.MapInfo.Tilemap.Bounds())
	}
	m.Canvas.Clear(color.Transparent)

	for index, l := range m.MapInfo.Tilemap.TileLayers {
		callback(m.Canvas, index)
		if l.Name == m.MapInfo.CollisionLayerName {
			//we do NOT render the collision layer
			continue
		}

		//if err := l.Draw(m.Canvas); err != nil {
		if err := m.lDraw(l, m.Canvas); err != nil {
			log.WithError(err).Error("GameMap.DrawAfter: could not draw layer")
			return err
		}
		if l.Name == m.MapInfo.HiddenLayer {
			m.MarkReRender = false
		}
	}

	for _, il := range m.MapInfo.Tilemap.ImageLayers {
		// The matrix shift is because images are drawn from the top-left in Tiled.
		if err := il.Draw(m.Canvas, pixel.IM.Moved(pixel.V(0, m.pixelHeight()))); err != nil {
			log.WithError(err).Error("Map.DrawAll: could not draw image layer")
			return err
		}
	}

	m.Canvas.Draw(target, mat.Moved(m.MapInfo.Tilemap.Bounds().Center()))

	return nil
}
func (m *GameMap) lDraw(l *tilepix.TileLayer, target *pixelgl.Canvas) error {
	// Only draw if the layer is dirty.
	if l.IsDirty || m.MarkReRender {
		// Initialise the batch
		if _, err := l.Batch(); err != nil {
			log.WithError(err).Error("TileLayer.Draw: could not get batch")
			return err
		}

		ts := l.Tileset
		numRows := ts.Tilecount / ts.Columns

		//(tileX)+(tileY * 100)
		// Loop through each decoded tile
		for tileIndex, tile := range l.DecodedTiles {
			if l.Name == m.MapInfo.HiddenLayer {
				if m.hideDecorationTile[tileIndex] {
					//do not render Hidden Layer
					continue
				}
			}
			// The Y component of the offset is set in Tiled from top down, setting here to negative because we want
			// from the bottom up.
			layerOffset := pixel.V(l.OffSetX, -l.OffSetY)
			tile.Draw(tileIndex, ts.Columns, numRows, ts, l.Batch_, layerOffset)
		}

		// Batch is drawn to, layer is no longer dirty.
		l.SetDirty(false)
	}

	l.Batch_.Draw(target)

	// Reset the dirty flag if the layer is not StaticBool
	if !l.StaticBool {
		l.SetDirty(true)
	}

	return nil
}

func (m GameMap) pixelHeight() float64 {
	return float64(m.MapInfo.Tilemap.Height * m.MapInfo.Tilemap.TileHeight)
}

func (m GameMap) GetTrigger(tileX, tileY float64) Trigger {
	x, y := m.GetTileIndex(tileX, tileY)
	return m.Triggers[[2]float64{x, y}]
}

func (m *GameMap) SetTrigger(tileX, tileY float64, t Trigger) {
	x, y := m.GetTileIndex(tileX, tileY)
	m.Triggers[[2]float64{x, y}] = t
}

//RemoveTrigger accept actual Tiled App coordinates e.g. 35, 22
func (m *GameMap) RemoveTrigger(tileX, tileY float64) {
	x, y := m.GetTileIndex(tileX, tileY)
	delete(m.Triggers, [2]float64{x, y})
}

//AddNPC helps in detecting player if x,y has NPC or not
func (m *GameMap) AddNPC(npc *Character) {
	m.NPCbyId[npc.Id] = npc
	m.NPCs = append(m.NPCs, npc)
	m.Entities = append(m.Entities, npc.Entity)
}

//WriteTile - accept actual Tiled App coordinates e.g. 35, 22
//it will bypassBlockedTile (removed collision)
func (m *GameMap) WriteTile(tileX, tileY float64, collision bool) {
	x, y := m.GetTileIndex(tileX, tileY)
	m.bypassBlockedTile[[2]float64{x, y}] = !collision
}
