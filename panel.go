package main

import (
	"github.com/faiface/pixel"
	"math"
)

type Panel struct {
	mTexture                pixel.Picture
	mTileSize, mCenterScale float64
	mUVs                    []pixel.Rect
	mTiles                  []*pixel.Sprite // the sprites representing the border.
	mSprites                []*pixel.Sprite
}

//PanelCreate
// {
//     texture = [texture],
//     size = [size of a single tile in pixels]
// }
func PanelCreate(texture pixel.Picture, size float64) Panel {

	p := Panel{
		mTexture:  texture,
		mUVs:      LoadAsFrames(texture, size, size),
		mTileSize: size,
	}

	// Fix up center U,Vs by moving them 0.5 texels in.
	p.mUVs[4] = PixelTexels(p.mUVs[4], p.mTexture)
	p.mCenterScale = p.mTileSize / (p.mTileSize - 1)

	// Create a sprite for each tile of the panel
	// 0. top left      1. top          2. top right
	// 3. left          4. middle       5. right
	// 6. bottom left   7. bottom       8. bottom right
	p.mTiles = make([]*pixel.Sprite, 9)
	for k, v := range p.mUVs {
		var sprite = pixel.NewSprite(texture, v)
		p.mTiles[k] = sprite
	}

	return p
}

func (p Panel) Position(bounds pixel.Rect) {
	// Reset scales
	//for _, v := range p.mTiles {
	//	v.Draw(global.gWin, pixel.IM.Scaled(pixel.V(0,0), 1))
	//}

	var hSize = p.mTileSize / 2
	// Align the corner tiles
	p.mTiles[6].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Min.X+hSize, bounds.Max.Y-hSize))) //topLeft
	p.mTiles[8].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Max.X-hSize, bounds.Max.Y-hSize))) //topRight
	p.mTiles[0].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Min.X+hSize, bounds.Min.Y+hSize))) //bottomLeft
	p.mTiles[2].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Max.X-hSize, bounds.Min.Y+hSize))) //bottomRight

	// Calculate how much to scale the side tiles
	var hWidth = bounds.W() / 2
	var widthScale float64 = math.Abs(hWidth-(2*p.mTileSize)) / (p.mTileSize / 2)
	var centerX = bounds.Center().X

	//top horizontal line
	p.mTiles[1].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, bounds.Max.Y-hSize)).ScaledXY(pixel.V(centerX, bounds.Max.Y-hSize), pixel.V(widthScale, 1)),
	)

	//bottom horizontal line
	p.mTiles[7].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, bounds.Min.Y+hSize)).ScaledXY(pixel.V(centerX, bounds.Min.Y+hSize), pixel.V(widthScale, 1)),
	)

	var hHeight = bounds.H() / 2
	var heightScale = math.Abs(hHeight-(2*p.mTileSize)) / (p.mTileSize / 2)
	var centerY = bounds.Center().Y

	//left vertical line
	p.mTiles[3].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(bounds.Min.X+hSize, centerY)).ScaledXY(pixel.V(bounds.Min.X+hSize, centerY), pixel.V(1, heightScale)),
	)

	//right vertical line
	p.mTiles[5].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(bounds.Max.X-hSize, centerY)).ScaledXY(pixel.V(bounds.Max.X-hSize, centerY), pixel.V(1, heightScale)),
	)

	// Scale the middle backing panel
	p.mTiles[4].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, centerY)).ScaledXY(
			pixel.V(centerX, centerY),
			pixel.V(widthScale, heightScale),
		),
	)

	// Hide corner tiles when scale is equal to zero
	//if left-right == 0 || top-bottom == 0 {
	//	for _, v := range p.mTiles {
	//		v.Draw(global.gWin, pixel.IM)
	//	}
	//}
}

func (p Panel) DrawAtPosition(v pixel.Vec, width, height float64) {
	p.Position(pixel.Rect{
		Min: pixel.V(v.X-width/2, v.Y-height/2),
		Max: pixel.V(v.X+width/2, v.Y+height/2),
	})
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
