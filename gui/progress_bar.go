package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/utilz"
)

type ProgressBar struct {
	x, y                             float64
	Background                       *pixel.Sprite
	Foreground                       *pixel.Sprite
	foregroundFrames                 []pixel.Rect
	foregroundFrame                  int
	scale                            float64
	foregroundPosition               pixel.Vec
	foregroundPng                    pixel.Picture
	Value, Maximum, HalfWidth, Width float64
}

func ProgressBarCreate(x, y float64, value, max float64, background, foreground string) ProgressBar {
	bgImg, err := utilz.LoadPicture(background)
	utilz.PanicIfErr(err)
	fgImg, err := utilz.LoadPicture(foreground)
	utilz.PanicIfErr(err)

	pb := ProgressBar{
		x:             x,
		y:             y,
		foregroundPng: fgImg,
		Background:    pixel.NewSprite(bgImg, bgImg.Bounds()),
		Foreground:    pixel.NewSprite(fgImg, fgImg.Bounds()),
		scale:         10,
		Value:         value,
		Maximum:       max,
	}

	// Get UV positions in texture atlas
	// A table with name fields: left, top, right, bottom
	pb.Width = pb.foregroundPng.Bounds().W()
	pb.HalfWidth = pb.Width / 2
	pb.foregroundFrames = utilz.LoadAsFrames(fgImg, pb.foregroundWidthBlock(), pb.foregroundPng.Bounds().H())

	pb.SetValue(value)

	return pb
}

func (pb ProgressBar) foregroundWidthBlock() float64 {
	return pb.Width * pb.scale / 100
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
	if pb.fallsInWhichPercent(pb.Value) < 10 {
		pb.Background.Draw(renderer, pixel.IM.Moved(mat))
	}

	fgMat := mat.Sub(pixel.V(pb.HalfWidth, 0))
	scaleFactor := pb.foregroundWidthBlock()
	for i := 0; i < pb.foregroundFrame; i++ {
		px := pixel.NewSprite(pb.foregroundPng, pb.foregroundFrames[i])
		px.Draw(renderer, pixel.IM.Moved(pixel.V(fgMat.X+(float64(i)*scaleFactor)+pb.scale, fgMat.Y)))
	}
}

/*
TO MATCH StackInterface below
*/
func (pb ProgressBar) HandleInput(win *pixelgl.Window) {
}
func (pb ProgressBar) Enter() {}
func (pb ProgressBar) Exit()  {}
func (pb ProgressBar) Update(dt float64) bool {
	return true
}
