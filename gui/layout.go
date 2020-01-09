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
	pos           pixel.Vec
	width, height float64
}

//centerPos would Player TileX, TileY
func LayoutCreate(x, y float64, win *pixelgl.Window) Layout {
	l := Layout{
		Panels: make(map[string]PanelDef, 0),
	}
	//first panel - full screen
	l.centerPos = pixel.V(x, y)
	fullScreenPanel := PanelDef{l.centerPos, win.Bounds().W() - 50, win.Bounds().H() - 50}
	l.Panels["screen"] = fullScreenPanel

	return l
}

func (l Layout) CreatePanel(name string) Panel {
	panelDef := l.getPanelDef(name)
	return PanelCreate(panelDef.pos, panelDef.width, panelDef.height)
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
	cam := pixel.IM.Scaled(l.centerPos, 1.0).Moved(l.centerPos)
	win.SetMatrix(cam)
}

//Contract will reduce the Panel dimension with give value
//e.g. Contract('screen', 118, 40)
func (l Layout) Contract(name string, horz, vert float64) {
	panelDef := l.getPanelDef(name)
	panelDef.width = panelDef.width - horz
	panelDef.height = panelDef.height - vert
}

//e.g. SplitHorz('screen', "top", "bottom", 0.12, 2) // x = from 0 to 1
func (l *Layout) SplitHorz(name, topName, bottomName string, x, margin float64) {
	parent := l.getPanelDef(name)
	//delete parent from Layout, we dont need it anymore
	delete(l.Panels, name)

	p1Height := parent.height * x
	p2Height := parent.height * (1 - x)

	l.Panels[topName] = PanelDef{
		pos:    pixel.V(parent.pos.X, parent.pos.Y-parent.height/2+p1Height/2-margin/2),
		width:  parent.width,
		height: p1Height - margin,
	}

	l.Panels[bottomName] = PanelDef{
		pos:    pixel.V(parent.pos.X, parent.pos.Y+parent.height/2-p2Height/2+margin/2),
		width:  parent.width,
		height: p2Height - margin,
	}
}

//e.g. SplitVert('bottom', "left", "party", 0.726, 2) // y = from 0 to 1
func (l *Layout) SplitVert(name, leftName, rightName string, y, margin float64) {
	parent := l.getPanelDef(name)
	//delete parent from Layout, we dont need it anymore
	delete(l.Panels, name)

	p1Width := parent.width * y
	p2Width := parent.width * (1 - y)

	l.Panels[rightName] = PanelDef{
		pos:    pixel.V(parent.pos.X+parent.width/2-p1Width/2+margin/2, parent.pos.Y),
		width:  p1Width - margin,
		height: parent.height,
	}

	l.Panels[leftName] = PanelDef{
		pos:    pixel.V(parent.pos.X-parent.width/2+p2Width/2-margin/2, parent.pos.Y),
		width:  p2Width - margin,
		height: parent.height,
	}
}

//since Panel renders from Center of X, y
func (l Layout) Top(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.Y + panel.height/2
}

func (l Layout) Bottom(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.Y - panel.height/2
}

func (l Layout) Left(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.X - panel.width/2
}

func (l Layout) Right(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.X + panel.width/2
}

func (l Layout) MidX(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.X
}

func (l Layout) MidY(name string) float64 {
	panel := l.getPanelDef(name)
	return panel.pos.Y
}
