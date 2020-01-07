package gui

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/globals"
)

type ProgressBar struct {
	x, y                      float64
	Background                *pixel.Sprite
	Foreground                *pixel.Sprite
	foregroundFrames          []pixel.Rect
	foregroundFrame           int
	scale                     float64
	foregroundPosition        pixel.Vec
	foregroundPng             pixel.Picture
	Value, Maximum, halfWidth float64
}

func ProgressBarCreate() ProgressBar {
	bgImg, fgImg := globals.ProgressBarBgPng, globals.ProgressBarFbPng
	pb := ProgressBar{
		x:             0,
		y:             0,
		foregroundPng: fgImg,
		Background:    pixel.NewSprite(bgImg, bgImg.Bounds()),
		Foreground:    pixel.NewSprite(fgImg, fgImg.Bounds()),
		scale:         10,
		Value:         0,
		Maximum:       100,
	}

	// Get UV positions in texture atlas
	// A table with name fields: left, top, right, bottom
	pb.halfWidth = bgImg.Bounds().W() / 2
	pb.foregroundFrames = globals.LoadAsFrames(fgImg, pb.foregroundWidthBlock(), pb.foregroundPng.Bounds().H())

	pb.SetValue(10)

	return pb
}

func (pb ProgressBar) foregroundWidthBlock() float64 {
	return pb.foregroundPng.Bounds().W() * pb.scale / 100
}

func (pb *ProgressBar) SetMax(maxHealth float64) {
	pb.Maximum = maxHealth
}
func (pb *ProgressBar) SetValue(health float64) {
	pb.Value = health
	pb.setForegroundValue()
}

func (pb *ProgressBar) setForegroundValue() {
	pb.foregroundFrame = pb.fallsInWhichPercent(pb.Value)
}

func (pb ProgressBar) fallsInWhichPercent(val float64) int {
	maxFramesToShow := val / pb.Maximum * 100

	if maxFramesToShow >= 90 {
		return 10
	}
	if maxFramesToShow >= 80 && maxFramesToShow < 90 {
		return 9
	}
	if maxFramesToShow >= 70 && maxFramesToShow < 80 {
		return 8
	}
	if maxFramesToShow >= 60 && maxFramesToShow < 70 {
		return 7
	}
	if maxFramesToShow >= 50 && maxFramesToShow < 60 {
		return 6
	}
	if maxFramesToShow >= 40 && maxFramesToShow < 50 {
		return 5
	}
	if maxFramesToShow >= 30 && maxFramesToShow < 40 {
		return 4
	}
	if maxFramesToShow >= 20 && maxFramesToShow < 30 {
		return 3
	}
	if maxFramesToShow >= 10 && maxFramesToShow < 20 {
		return 2
	}
	if maxFramesToShow >= 0 && maxFramesToShow < 10 {
		return 1
	}
	return 0
}

func (pb *ProgressBar) SetPosition(x, y float64) {
	pb.x = x
	pb.y = y
}
func (pb ProgressBar) GetPosition() (x, y float64) {
	return pb.x, pb.y
}

func (pb ProgressBar) Render(renderer pixel.Target) {
	mat := pixel.V(pb.x, pb.y)
	pb.Background.Draw(renderer, pixel.IM.Moved(mat))

	fgMat := mat.Sub(pixel.V(pb.halfWidth, 0))
	scaleFactor := pb.foregroundWidthBlock()
	for i := 0; i < pb.foregroundFrame; i++ {
		px := pixel.NewSprite(pb.foregroundPng, pb.foregroundFrames[i])
		px.Draw(renderer, pixel.IM.Moved(
			pixel.V(fgMat.X+(float64(i)*scaleFactor)+pb.scale, fgMat.Y)))
	}
}
