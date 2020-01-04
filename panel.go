package main

import (
	"github.com/faiface/pixel"
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
	//var center = p.mUVs[4]
	//var pixelToTexelX = 1 / p.mTexture.Bounds().W()
	//var pixelToTexelY = 1 / p.mTexture.Bounds().H()
	//center.Min.X = center.Min.X + (pixelToTexelX / 2)
	//center.Min.Y = center.Min.Y + (pixelToTexelY / 2)
	//center.Max.X = center.Max.X - (pixelToTexelX / 2)
	//center.Max.Y = center.Max.Y - (pixelToTexelY / 2)
	//p.mUVs[4] = center

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

	//var hSize = p.mTileSize/2
	// Align the corner tiles
	p.mTiles[6].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Min.X, bounds.Max.Y))) //topLeft
	p.mTiles[8].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Max.X, bounds.Max.Y))) //topRight
	p.mTiles[0].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Min.X, bounds.Min.Y))) //bottomLeft
	p.mTiles[2].Draw(global.gWin, pixel.IM.Moved(pixel.V(bounds.Max.X, bounds.Min.Y))) //bottomRight

	// Calculate how much to scale the side tiles
	var percent float64 = 36
	var hWidth = bounds.W() / 2
	var widthScale = hWidth - (hWidth * percent / 100)
	var centerX = bounds.Center().X

	//top horizontal line
	p.mTiles[1].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, bounds.Max.Y)).ScaledXY(pixel.V(centerX, bounds.Max.Y), pixel.V(widthScale, 1)),
	)

	//bottom horizontal line
	p.mTiles[7].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, bounds.Min.Y)).ScaledXY(pixel.V(centerX, bounds.Min.Y), pixel.V(widthScale, 1)),
	)

	var hHeight = bounds.H() / 2
	var heightScale = hHeight - (hHeight * percent / 100)
	var centerY = bounds.Center().Y

	//left vertical line
	p.mTiles[3].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(bounds.Min.X, centerY)).ScaledXY(pixel.V(bounds.Min.X, centerY), pixel.V(1, heightScale)),
	)

	//right vertical line
	p.mTiles[5].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(bounds.Max.X, centerY)).ScaledXY(pixel.V(bounds.Max.X, centerY), pixel.V(1, heightScale)),
	)

	// Scale the middle backing panel
	p.mTiles[4].Draw(
		global.gWin,
		pixel.IM.Moved(pixel.V(centerX, centerY)).ScaledXY(
			pixel.V(centerX, centerY),
			pixel.V(widthScale, heightScale)),
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
