package main

import (
	"github.com/faiface/pixel"
	"github.com/salviati/go-tmx/tmx"
)

type GameMap struct {
	mX, mY int

	// To track the camera position
	mCamX, mCamY int

	mTilemap *tmx.Map
	mSprites map[string]*pixel.Sprite

	mTileSprite pixel.Sprite
	mWidth, mHeight int

	mTiles []*pixel.Batch
	mTilesIndices map[string]int
	mTilesCounter int

	mTileWidth, mTileHeight int
}


func (m *GameMap) Create(tilemap *tmx.Map) {
	// assuming exported tiled map
	//lua definition has 1 layer
	m.mTilemap = tilemap

	//Top left corner of the map
	m.mX = -global.gWindowWidth / 2
	m.mY = global.gWindowHeight / 2

	m.mHeight = tilemap.Height
	m.mWidth = tilemap.Width

	m.mTileWidth = tilemap.TileWidth
	m.mTileHeight = tilemap.TileHeight

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

func getTileLocation(tID int, numColumns int, numRows int) (x, y int) {
	x = tID % numColumns
	y = numRows - (tID / numColumns) - 1
	return
}

func (m GameMap) getTilePos(idx int) pixel.Vec {
	width := m.mTilemap.Width
	height := m.mTilemap.Height
	gamePos := pixel.V(
		float64(idx % width) - 1,
		float64(height)-float64(idx / width),
	)
	return gamePos
}

func (m GameMap) Render() {
	// Draw tiles
	for _, batch := range m.mTiles {
		batch.Clear()
	}

	for _, layer := range m.mTilemap.Layers {
		for tileIndex, tile := range layer.DecodedTiles {
			ts := layer.Tileset
			tID := int(tile.ID)

			if tID == 0 {
				// Tile ID 0 means blank, skip it.
				continue
			}

			// Calculate the framing for the tile within its tileset's source image
			numRows := ts.Tilecount / ts.Columns
			x, y := getTileLocation(tID, ts.Columns, numRows)
			gamePos := m.getTilePos(tileIndex)

			iX := float64(x) * float64(ts.TileWidth)
			fX := iX + float64(ts.TileWidth)
			iY := float64(y) * float64(ts.TileHeight)
			fY := iY + float64(ts.TileHeight)

			sprite := m.mSprites[ts.Image.Source]
			sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
			sprite.Draw(m.mTiles[m.mTilesIndices[ts.Image.Source]], pixel.IM.Moved(pos))
		}
	}

	for _, batch := range m.mTiles {
		batch.Draw(global.gWin)
	}
}
