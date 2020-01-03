package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"image/color"
)

func DrawFixedTopPanel(gMap *GameMap) {
	//Panel
	gWidth, gHeight := global.gWindowWidth/2, global.gWindowHeight/2
	panelTopLeftX, panelTopLeftY := gMap.mCamX-gWidth+10, gMap.mCamY+gHeight-10
	panelBottomRightX, panelBottomRightY := gMap.mCamX+gWidth-10, gMap.mCamY+150

	panel.Clear()
	panel.Color = color.Black
	panel.Push(
		pixel.V(panelTopLeftX, panelTopLeftY),
		pixel.V(panelBottomRightX, panelBottomRightY),
	)
	panel.Rectangle(0)
	panel.Draw(global.gWin)

	DrawText(panelTopLeftX, panelTopLeftY, `
A nation can survive its fools, and even the ambitious. But 
it cannot survive treason from within. An enemy at the gates
is less formidable, for he is known and carries his banner openly.
`)
}

func DrawText(panelTopLeftX, panelTopLeftY float64, textMsg string) {
	//Text
	//var textAddToCenter float64 = panelBottomRightX/2+10
	var paddingX, paddingY float64 = 10, -5
	textPos := pixel.V(panelTopLeftX+paddingX, panelTopLeftY+paddingY)
	basicTxt := text.New(textPos, basicAtlas)
	basicTxt.LineHeight = basicAtlas.LineHeight() * 1.5
	fmt.Fprintln(basicTxt, textMsg)
	basicTxt.Draw(global.gWin, pixel.IM)
}
