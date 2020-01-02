package main

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	log "github.com/sirupsen/logrus"
	"image/color"
)

type GameMap struct {
	mX, mY float64

	// To track the camera position
	mCamX, mCamY float64

	mTilemap *tilepix.Map
	mSprites map[string]*pixel.Sprite

	mTileSprite     pixel.Sprite
	mWidth, mHeight float64

	mTiles        []*pixel.Batch
	mTilesIndices map[string]int
	mTilesCounter int

	mTileWidth, mTileHeight float64
	blockingTileGID         tilepix.GID
	canvas                  *pixelgl.Canvas
	renderLayer             int

	mTriggers map[[2]float64]Trigger
}

func (m *GameMap) Create(tilemap *tilepix.Map) {
	// assuming exported tiled map
	//lua definition has 1 layer
	m.mTilemap = tilemap
	m.mTriggers = make(map[[2]float64]Trigger)

	m.mHeight = float64(tilemap.Height)
	m.mWidth = float64(tilemap.Width)

	m.mTileWidth = float64(tilemap.TileWidth)
	m.mTileHeight = float64(tilemap.TileHeight)

	//Bottom left corner of the map, since pixel starts at 0, 0
	m.mX = m.mTileWidth
	m.mY = m.mTileHeight

	m.canvas = pixelgl.NewCanvas(m.mTilemap.Bounds())
	m.setTiles()
	m.setBlockingTileInfo()
}

func (m *GameMap) setBlockingTileInfo() {
	for _, tile := range m.mTilemap.Tilesets {
		if tile.Name == "collision_px" {
			m.blockingTileGID = tile.FirstGID
			break
		}
	}
}

//IsBlockingTile check's x, y cords on collision map layer
// if ID is not 0, tile exists on x, y we return true
func (m GameMap) IsBlockingTile(x, y, layer int) bool {
	tile := m.mTilemap.TileLayers[layer].DecodedTiles[x+y*int(m.mWidth)]
	return !tile.IsNil() || tile.ID != 0
}

func (m *GameMap) setTiles() {
	batches := make([]*pixel.Batch, 0)
	batchIndices := make(map[string]int)
	batchCounter := 0

	// Load the sprites
	sprites := make(map[string]*pixel.Sprite)
	for _, tileset := range m.mTilemap.Tilesets {
		if _, alreadyLoaded := sprites[tileset.Image.Source]; !alreadyLoaded {
			sprite, pictureData := LoadSprite(tileset.Image.Source)
			sprites[tileset.Image.Source] = sprite
			batches = append(batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
			batchIndices[tileset.Image.Source] = batchCounter
			batchCounter++
		}
	}
	m.mTiles = batches
	m.mTilesIndices = batchIndices
	m.mTilesCounter = batchCounter
	m.mSprites = sprites
}

//CamToTile pan camera to given coordinates
func (m *GameMap) CamToTile(x, y float64) {
	tileX, tileY := m.GetTileIndex(x, y)
	x = tileX - m.mTileWidth/2
	y = tileY - m.mTileHeight/2
	m.Goto(x, y)
}

func (m *GameMap) Goto(x, y float64) {
	m.mCamX = x
	m.mCamY = y
}

func (m GameMap) GetTileIndex(x, y float64) (tileX, tileY float64) {
	y = m.mHeight - y //make count y from top (Tiled app starts from top)
	tileX = m.mX + (x * m.mTileWidth)
	tileY = m.mY + (y * m.mTileHeight)
	return
}

func (m GameMap) GetTilePositionAtFeet(x, y, charW, charH float64) pixel.Vec {
	tileX, tileY := m.GetTileIndex(x, y)
	x = tileX - charW/2
	y = tileY - charH/2
	return pixel.V(x, y)
}

func (m GameMap) DrawAll(target pixel.Target, clearColour color.Color, mat pixel.Matrix) {
	//m.mTilemap.DrawAll(global.gWin, color.Transparent, pixel.IM)
	m.mTilemap.DrawAll(target, clearColour, mat)
}

//DrawAfter will render the callback function after given layer index
// uses pixelgl Canvas instead of gWin to render
func (m GameMap) DrawAfter(layer int, callback func(canvas *pixelgl.Canvas)) error {
	// Draw tiles
	target, mat := global.gWin, pixel.IM

	if m.canvas == nil {
		m.canvas = pixelgl.NewCanvas(m.mTilemap.Bounds())
	}
	m.canvas.Clear(color.Transparent)

	for index, l := range m.mTilemap.TileLayers {
		//we do NOT render the collision layer
		if l.Name == "collision" {
			continue
		}
		if index == layer {
			callback(m.canvas)
		}
		if err := l.Draw(m.canvas); err != nil {
			log.WithError(err).Error("Map.DrawAll: could not draw layer")
			return err
		}
	}

	for _, il := range m.mTilemap.ImageLayers {
		// The matrix shift is because images are drawn from the top-left in Tiled.
		if err := il.Draw(m.canvas, pixel.IM.Moved(pixel.V(0, m.pixelHeight()))); err != nil {
			log.WithError(err).Error("Map.DrawAll: could not draw image layer")
			return err
		}
	}

	m.canvas.Draw(target, mat.Moved(m.mTilemap.Bounds().Center()))

	return nil
}

func (m GameMap) pixelHeight() float64 {
	return float64(m.mTilemap.Height * m.mTilemap.TileHeight)
}

func (m GameMap) GetTrigger(x, y float64) Trigger {
	return m.mTriggers[[2]float64{x, y}]
}
func (m GameMap) SetTrigger(x, y float64, t Trigger) {
	m.mTriggers[[2]float64{x, y}] = t
}
