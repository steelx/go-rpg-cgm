package gui

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/utilz"
	"math"
)

type Panel struct {
	mBounds                 pixel.Rect
	mTexture                pixel.Picture
	mTileSize, mCenterScale float64
	mUVs                    []pixel.Rect
	mTiles                  []*pixel.Sprite // the sprites representing the border.
	mSprites                []*pixel.Sprite
}

//PanelCreate
func PanelCreate(pos pixel.Vec, width, height float64) Panel {
	var size float64 = 3
	panelPng, err := utilz.LoadPicture("../resources/simple_panel.png")
	utilz.PanicIfErr(err)
	p := Panel{
		mTexture:  panelPng,
		mUVs:      utilz.LoadAsFrames(panelPng, size, size),
		mTileSize: size,
		mBounds: pixel.Rect{
			Min: pixel.V(pos.X-width/2, pos.Y-height/2),
			Max: pixel.V(pos.X+width/2, pos.Y+height/2),
		},
	}

	p.RefreshPanelCorners()

	return p
}

func (p *Panel) RefreshPanelCorners() {
	// Fix up center U,Vs by moving them 0.5 texels in.
	p.mUVs[4] = PixelTexels(p.mUVs[4], p.mTexture)
	p.mCenterScale = p.mTileSize / (p.mTileSize - 1)

	// Create a sprite for each tile of the panel
	// 0. top left      1. top          2. top right
	// 3. left          4. middle       5. right
	// 6. bottom left   7. bottom       8. bottom right
	p.mTiles = make([]*pixel.Sprite, 9)
	for k, v := range p.mUVs {
		var sprite = pixel.NewSprite(p.mTexture, v)
		p.mTiles[k] = sprite
	}
}

func (p Panel) GetCorners() (topLeft pixel.Vec, topRight pixel.Vec, bottomLeft pixel.Vec, bottomRight pixel.Vec) {
	var hSize = p.mTileSize / 2
	topLeft = pixel.V(p.mBounds.Min.X+hSize, p.mBounds.Max.Y-hSize)
	topRight = pixel.V(p.mBounds.Max.X-hSize, p.mBounds.Max.Y-hSize)
	bottomLeft = pixel.V(p.mBounds.Min.X+hSize, p.mBounds.Min.Y+hSize)
	bottomRight = pixel.V(p.mBounds.Max.X-hSize, p.mBounds.Min.Y+hSize)
	return
}

func (p Panel) Draw(renderer pixel.Target) {
	bounds := p.mBounds

	// Align the corner tiles
	topLeft, topRight, bottomLeft, bottomRight := p.GetCorners()
	p.mTiles[6].Draw(renderer, pixel.IM.Moved(topLeft))
	p.mTiles[8].Draw(renderer, pixel.IM.Moved(topRight))
	p.mTiles[0].Draw(renderer, pixel.IM.Moved(bottomLeft))
	p.mTiles[2].Draw(renderer, pixel.IM.Moved(bottomRight))

	// Calculate how much to scale the side tiles
	var hSize = p.mTileSize / 2
	var hWidth = bounds.W() / 2
	var widthScale float64 = math.Abs(hWidth-(2*p.mTileSize)) / (p.mTileSize / 2)
	var centerX = bounds.Center().X

	//top horizontal line
	p.mTiles[1].Draw(
		renderer,
		pixel.IM.Moved(pixel.V(centerX, bounds.Max.Y-hSize)).ScaledXY(pixel.V(centerX, bounds.Max.Y-hSize), pixel.V(widthScale, 1)),
	)

	//bottom horizontal line
	p.mTiles[7].Draw(
		renderer,
		pixel.IM.Moved(pixel.V(centerX, bounds.Min.Y+hSize)).ScaledXY(pixel.V(centerX, bounds.Min.Y+hSize), pixel.V(widthScale, 1)),
	)

	var hHeight = bounds.H() / 2
	var heightScale = math.Abs(hHeight-(2*p.mTileSize)) / (p.mTileSize / 2)
	var centerY = bounds.Center().Y

	//left vertical line
	p.mTiles[3].Draw(
		renderer,
		pixel.IM.Moved(pixel.V(bounds.Min.X+hSize, centerY)).ScaledXY(pixel.V(bounds.Min.X+hSize, centerY), pixel.V(1, heightScale)),
	)

	//right vertical line
	p.mTiles[5].Draw(
		renderer,
		pixel.IM.Moved(pixel.V(bounds.Max.X-hSize, centerY)).ScaledXY(pixel.V(bounds.Max.X-hSize, centerY), pixel.V(1, heightScale)),
	)

	// Scale the middle backing panel
	p.mTiles[4].Draw(
		renderer,
		pixel.IM.Moved(pixel.V(centerX, centerY)).ScaledXY(
			pixel.V(centerX, centerY),
			pixel.V(widthScale+(hSize*p.mCenterScale)+p.mTileSize, heightScale+(hSize*p.mCenterScale)),
		),
	)
}

func PixelTexels(pix pixel.Rect, texture pixel.Picture) pixel.Rect {
	var pixelToTexelX = 1 / texture.Bounds().W()
	var pixelToTexelY = 1 / texture.Bounds().H()
	pix.Min.X = pix.Min.X + (pixelToTexelX / 2)
	pix.Min.Y = pix.Min.Y + (pixelToTexelY / 2)
	pix.Max.X = pix.Max.X - (pixelToTexelX / 2)
	pix.Max.Y = pix.Max.Y - (pixelToTexelY / 2)
	return pix
}
