package main

import (
	"github.com/faiface/pixel"
)

type CharacterDefinition struct {
	texture       pixel.Picture
	width, height float64
	startFrame    int
	tileX, tileY  float64
}

//Entity represents any kind of map object from a
// treasure chest to an NPC
type Entity struct {
	mSprite         *pixel.Sprite
	mTexture        pixel.Batch
	mHeight, mWidth float64
	mTileX, mTileY  float64
	startFrame      int
	mFrames         []pixel.Rect
}

func CreateEntity(def CharacterDefinition) *Entity {
	e := &Entity{}

	e.mFrames = LoadAsFrames(def.texture, def.width, def.height)
	e.mSprite = pixel.NewSprite(def.texture, e.mFrames[def.startFrame])
	e.mWidth = def.width
	e.mHeight = def.height
	e.mTileX = def.tileX
	e.mTileY = def.tileY
	e.startFrame = def.startFrame
	return e
}

func (e *Entity) SetFrame(frame int) {
	e.startFrame = frame
}

//TeleportAndDraw hero movement & set position for sprite
func (e *Entity) TeleportAndDraw(gMap GameMap) {
	//spriteFrame := e.mFrames[e.startFrame]
	vec := gMap.GetTilePositionAtFeet(e.mTileX, e.mTileY, e.mWidth, e.mHeight)
	e.mSprite.Draw(global.gWin, pixel.IM.Moved(vec))
}
