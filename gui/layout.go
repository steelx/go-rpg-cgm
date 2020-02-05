package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"log"
)

type Layout struct {
	centerPos pixel.Vec
	Panels    map[string]PanelDef
}

type PanelDef struct {
	Pos           pixel.Vec
	Width, Height float64
}

func (p PanelDef) GetSize() (width, height float64) {
	return p.Width, p.Height
}

//centerPos would Player TileX, TileY
func LayoutCreate(x, y float64, win *pixelgl.Window) Layout {
	l := Layout{
		Panels: make(map[string]PanelDef, 0),
	}
	//first panel - full screen
	l.centerPos = pixel.V(x, y)
	var width, height float64
	if win.Monitor() != nil {
		width, height = win.Monitor().Size()
	} else {
		width, height = win.Bounds().W(), win.Bounds().H()
	}
	fullScreenPanel := PanelDef{l.centerPos, width - 50, height - 50}
	l.Panels["screen"] = fullScreenPanel

	return l
}

func (l Layout) CreatePanel(name string) Panel {
	panelDef := l.getPanelDef(name)
	return PanelCreate(panelDef.Pos, panelDef.Width, panelDef.Height)
}
func (l Layout) getPanelDef(name string) PanelDef {
	panelDef, ok := l.Panels[name]
	if !ok {
		log.Fatal("Layout " + name + " doesnt exist!")
	}
	return panelDef
}

//DebugRender just for debug
func (l Layout) DebugRender(win *pixelgl.Window) {
	for key := range l.Panels {
		panel := l.CreatePanel(key)
		panel.Draw(win)
	}

	//temp camera matrix
	//cam := pixel.IM.Scaled(l.centerPos, 1.0).Moved(l.centerPos)
	//win.SetMatrix(cam)
}

//Contract will reduce the Panel dimension with give value
//e.g. Contract('screen', 118, 40)
func (l Layout) Contract(name string, horz, vert float64) {
	panelDef := l.getPanelDef(name)
	panelDef.Width = panelDef.Width - horz
	panelDef.Height = panelDef.Height - vert
}

//e.g. SplitHorz('screen', "top", "bottom", 0.12, 2) // X = from 0 to 1
// X represents top height percent
func (l *Layout) SplitHorz(name, topName, bottomName string, x, margin float64) {
	parent := l.getPanelDef(name)
	//delete parent from Layout, we dont need it anymore
	delete(l.Panels, name)

	p1Height := parent.Height * x
	p2Height := parent.Height * (1 - x)

	l.Panels[topName] = PanelDef{
		Pos:    pixel.V(parent.Pos.X, parent.Pos.Y+parent.Height/2-p1Height/2+margin/2),
		Width:  parent.Width,
		Height: p1Height - margin,
	}

	l.Panels[bottomName] = PanelDef{
		Pos:    pixel.V(parent.Pos.X, parent.Pos.Y-parent.Height/2+p2Height/2-margin/2),
		Width:  parent.Width,
		Height: p2Height - margin,
	}
}

//e.g. SplitVert('bottom', "left", "party", 0.726, 2) // Y = from 0 to 1
// Y represents rightName width percent
func (l *Layout) SplitVert(name, leftName, rightName string, y, margin float64) {
	parent := l.getPanelDef(name)
	//delete parent from Layout, we dont need it anymore
	delete(l.Panels, name)

	p1Width := parent.Width * y
	p2Width := parent.Width * (1 - y)

	l.Panels[rightName] = PanelDef{
		Pos:    pixel.V(parent.Pos.X+parent.Width/2-p1Width/2+margin/2, parent.Pos.Y),
		Width:  p1Width - margin,
		Height: parent.Height,
	}

	l.Panels[leftName] = PanelDef{
		Pos:    pixel.V(parent.Pos.X-parent.Width/2+p2Width/2-margin/2, parent.Pos.Y),
		Width:  p2Width - margin,
		Height: parent.Height,
	}
}

//since Panel renders from Center of X, Y
func (l Layout) Top(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.Y + panel.Height/2
}

func (l Layout) Bottom(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.Y - panel.Height/2
}

func (l Layout) Left(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.X - panel.Width/2
}

func (l Layout) Right(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.X + panel.Width/2
}

func (l Layout) MidX(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.X
}

func (l Layout) MidY(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.Pos.Y
}
