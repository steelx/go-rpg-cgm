package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
)

type Panel struct {
	mBounds                 pixel.Rect
	mTileSize, mCenterScale float64
	BGColor                 color.RGBA //hex
}

var imd = imdraw.New(nil)

//PanelCreate
func PanelCreate(pos pixel.Vec, width, height float64) Panel {
	var size float64 = 3
	p := Panel{
		mTileSize: size,
		mBounds: pixel.Rect{
			Min: pixel.V(pos.X-width/2, pos.Y-height/2),
			Max: pixel.V(pos.X+width/2, pos.Y+height/2),
		},
		BGColor: utilz.HexToColor("#002D64"),
	}

	return p
}

func (p Panel) GetCorners() (topLeft pixel.Vec, topRight pixel.Vec, bottomLeft pixel.Vec, bottomRight pixel.Vec) {
	var hSize = p.mTileSize / 2
	topLeft = pixel.V(p.mBounds.Min.X+hSize, p.mBounds.Max.Y-hSize)
	topRight = pixel.V(p.mBounds.Max.X-hSize, p.mBounds.Max.Y-hSize)
	bottomLeft = pixel.V(p.mBounds.Min.X+hSize, p.mBounds.Min.Y+hSize)
	bottomRight = pixel.V(p.mBounds.Max.X-hSize, p.mBounds.Min.Y+hSize)
	return
}

func (p *Panel) Draw(renderer pixel.Target) {
	imd.Clear()
	topLeft, topRight, bottomLeft, bottomRight := p.GetCorners()

	// Middle backing panel
	imd.Color = p.BGColor
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(
		bottomLeft.Add(pixel.V(1, 1)),
		topLeft.Add(pixel.V(1, -1)),
		bottomRight.Add(pixel.V(-1, 1)),
		topRight.Add(pixel.V(-1, -1)),
	)
	imd.Rectangle(0)

	// Align the corner tiles, Rectangle above
	imd.Color = utilz.HexToColor("#D9AO66")
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(topLeft)
	imd.Push(topRight)
	imd.Push(bottomLeft)
	imd.Push(bottomRight)
	imd.Rectangle(1)

	imd.Draw(renderer)
}
