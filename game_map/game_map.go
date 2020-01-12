package game_map

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	log "github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/globals"
	"image/color"
)

type GameMap struct {
	x, y float64

	// To track the camera position
	CamX, CamY float64

	Tilemap            *tilepix.Map
	CollisionLayer     int
	CollisionLayerName string
	sprites            map[string]*pixel.Sprite

	mTileSprite   pixel.Sprite
	Width, Height float64

	Tiles        []*pixel.Batch
	tilesIndices map[string]int
	tilesCounter int

	TileWidth, TileHeight float64
	blockingTileGID       tilepix.GID
	Canvas                *pixelgl.Canvas
	renderLayer           int

	Triggers map[[2]float64]Trigger
	Entities []*Entity
	NPCs     []*Character
}

func MapCreate(tilemap *tilepix.Map, collisionLayer int, collisionLayerName string) *GameMap {
	m := &GameMap{
		Tilemap:            tilemap,
		CollisionLayer:     collisionLayer,
		CollisionLayerName: collisionLayerName,
	}

	m.Triggers = make(map[[2]float64]Trigger)
	m.Entities = make([]*Entity, 0)

	m.Height = float64(tilemap.Height)
	m.Width = float64(tilemap.Width)

	m.TileWidth = float64(tilemap.TileWidth)
	m.TileHeight = float64(tilemap.TileHeight)

	//Bottom left corner of the map, since pixel starts at 0, 0
	m.x = m.TileWidth
	m.y = m.TileHeight

	m.Canvas = pixelgl.NewCanvas(m.Tilemap.Bounds())
	m.setTiles()
	m.setBlockingTileInfo()
	return m
}

func (m *GameMap) setBlockingTileInfo() {
	for _, tile := range m.Tilemap.Tilesets {
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
func (m GameMap) IsBlockingTile(x, y int) bool {
	tile := m.Tilemap.TileLayers[m.CollisionLayer].DecodedTiles[x+y*int(m.Width)]
	return !tile.IsNil() || tile.ID != 0
}

func (m *GameMap) setTiles() {
	batches := make([]*pixel.Batch, 0)
	batchIndices := make(map[string]int)
	batchCounter := 0

	// Load the sprites
	sprites := make(map[string]*pixel.Sprite)
	for _, tileset := range m.Tilemap.Tilesets {
		if _, alreadyLoaded := sprites[tileset.Image.Source]; !alreadyLoaded {
			sprite, pictureData := globals.LoadSprite(tileset.Image.Source)
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

func (m GameMap) GetTileIndex(x, y float64) (tileX, tileY float64) {
	y = m.Height - y //make count Y from top (Tiled app starts from top)
	tileX = m.x + (x * m.TileWidth)
	tileY = m.y + (y * m.TileHeight)
	return
}

func (m GameMap) GetTilePositionAtFeet(x, y, charW, charH float64) pixel.Vec {
	tileX, tileY := m.GetTileIndex(x, y)
	x = tileX - charW/2
	y = tileY - charH/2
	return pixel.V(x, y)
}

func (m GameMap) DrawAll(target pixel.Target, clearColour color.Color, mat pixel.Matrix) {
	//m.Tilemap.DrawAll(Global.Win, color.Transparent, pixel.IM)
	m.Tilemap.DrawAll(target, clearColour, mat)
}

//DrawAfter will render the callback function after given layer index
// uses pixelgl Canvas instead of Win to render
func (m GameMap) DrawAfter(callback func(canvas *pixelgl.Canvas, layer int)) error {
	// Draw tiles
	target, mat := globals.Global.Win, pixel.IM

	if m.Canvas == nil {
		m.Canvas = pixelgl.NewCanvas(m.Tilemap.Bounds())
	}
	m.Canvas.Clear(color.Transparent)

	for index, l := range m.Tilemap.TileLayers {
		callback(m.Canvas, index)
		if l.Name == m.CollisionLayerName {
			//we do NOT render the collision layer
			continue
		}
		if err := l.Draw(m.Canvas); err != nil {
			log.WithError(err).Error("GameMap.DrawAfter: could not draw layer")
			return err
		}
	}

	for _, il := range m.Tilemap.ImageLayers {
		// The matrix shift is because images are drawn from the top-left in Tiled.
		if err := il.Draw(m.Canvas, pixel.IM.Moved(pixel.V(0, m.pixelHeight()))); err != nil {
			log.WithError(err).Error("Map.DrawAll: could not draw image layer")
			return err
		}
	}

	m.Canvas.Draw(target, mat.Moved(m.Tilemap.Bounds().Center()))

	return nil
}

func (m GameMap) pixelHeight() float64 {
	return float64(m.Tilemap.Height * m.Tilemap.TileHeight)
}

func (m GameMap) GetTrigger(x, y float64) Trigger {
	return m.Triggers[[2]float64{x, y}]
}
func (m GameMap) SetTrigger(x, y float64, t Trigger) {
	tileX, tileY := m.GetTileIndex(x, y)
	m.Triggers[[2]float64{tileX, tileY}] = t
}

//AddNPC helps in detecting player if x,y has NPC or not
func (m *GameMap) AddNPC(npc *Character) {
	m.NPCs = append(m.NPCs, npc)
	m.Entities = append(m.Entities, npc.Entity)
}
