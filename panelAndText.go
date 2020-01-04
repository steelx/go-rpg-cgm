package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
)

func DrawPanelFixedTop(gMap *GameMap, textMsg string, textAtlas *text.Atlas) *imdraw.IMDraw {
	//Panel
	gWidth, gHeight := global.gWindowWidth/2, global.gWindowHeight/2
	panelTopLeftX, panelTopLeftY := gMap.mCamX-gWidth+10, gMap.mCamY+gHeight-10
	panelBottomRightX, panelBottomRightY := gMap.mCamX+gWidth-10, gMap.mCamY+150

	panel := imdraw.New(nil)
	//panel.Clear()
	panel.Color = color.Black
	panel.Push(
		pixel.V(panelTopLeftX, panelTopLeftY),
		pixel.V(panelBottomRightX, panelBottomRightY),
	)
	panel.Rectangle(0)
	panel.Draw(global.gWin)

	return panel
}

func DrawText(pos pixel.Vec, textMsg string, textAtlas *text.Atlas) {
	//Text
	textPos := pos
	basicTxt := text.New(textPos, textAtlas)
	basicTxt.LineHeight = textAtlas.LineHeight() * 1.5
	fmt.Fprintln(basicTxt, textMsg)
	basicTxt.Draw(global.gWin, pixel.IM)
}

func DrawPanelCharacterTop(player *Entity, textMsg string, textAtlas *text.Atlas) *imdraw.IMDraw {
	//Panel
	var gWidth, gHeight float64 = 150, 70
	x, y := player.gMap.GetTileIndex(player.mTileX, player.mTileY)
	panelTopLeftX, panelTopLeftY := x-gWidth-10, y+gHeight+10
	panelBottomRightX, panelBottomRightY := x+gWidth-10, y+20

	panel := imdraw.New(nil)
	//panel.Clear()
	panel.Color = colornames.Darkviolet
	panel.Push(
		pixel.V(panelTopLeftX, panelTopLeftY),
		pixel.V(panelBottomRightX, panelBottomRightY),
	)
	panel.Rectangle(0)
	panel.Draw(global.gWin)
	DrawText(pixel.V(x, y), textMsg, textAtlas)
	return panel
}
