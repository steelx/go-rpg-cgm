package main

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
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
}

func (m *GameMap) Create(tilemap *tilepix.Map) {
	// assuming exported tiled map
	//lua definition has 1 layer
	m.mTilemap = tilemap

	m.mHeight = float64(tilemap.Height)
	m.mWidth = float64(tilemap.Width)

	m.mTileWidth = float64(tilemap.TileWidth)
	m.mTileHeight = float64(tilemap.TileHeight)

	//Bottom left corner of the map, since pixel starts at 0, 0
	m.mX = m.mTileWidth
	m.mY = m.mTileHeight

	m.SetTiles()
}

func (m *GameMap) SetTiles() {
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

func (m *GameMap) CamToTile(x, y float64) {
	y = m.mHeight - y
	x = x - 1

	x = m.mX + (x * m.mTileWidth) - m.mTileWidth/2
	y = m.mY + (y * m.mTileHeight) - m.mTileHeight/2
	m.Goto(x, y)
}

func (m *GameMap) Goto(x, y float64) {
	m.mCamX = x
	m.mCamY = y
}

func (m *GameMap) GetTilePositionAtFeet(x, y, charW, charH float64) pixel.Vec {
	y = m.mHeight - y //make count y from top (Tiled app starts from top)
	//x = x - 1
	x = m.mX + (x * m.mTileWidth) - charW/2
	y = m.mY + (y * m.mTileHeight) - charH/2

	return pixel.V(x, y)
}

func getTileLocation(tID, numColumns, numRows int) (x, y int) {
	x = tID % numColumns
	y = numRows - (tID / numColumns) - 1
	return
}

func (m GameMap) getTilePos(idx int) pixel.Vec {
	width := m.mTilemap.Width
	height := m.mTilemap.Height
	gamePos := pixel.V(
		float64(idx%width)-1,
		float64(height)-float64(idx/width),
	)
	return gamePos
}

func (m GameMap) Render() {
	// Draw tiles
	m.mTilemap.DrawAll(global.gWin, color.Transparent, pixel.IM)
}
