package main

import (
	"github.com/faiface/pixel"
)

type CharacterDefinition struct {
	texture       pixel.Picture
	width, height int
	startFrame    int
	tileX, tileY  int
}

//Entity represents any kind of map object from a
// treasure chest to an NPC
type Entity struct {
	mSprite            *pixel.Sprite
	mTexture           pixel.Batch
	mHeight, mWidth    int
	mTileX, mTileY     int
	startFrame, mFrame int
	mFrames            []pixel.Rect
}

func CreateEntity(def CharacterDefinition) *Entity {
	e := &Entity{}

	e.mFrames = LoadAsFrames(def.texture, float64(def.width), float64(def.height))
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
	spriteFrame := e.mFrames[e.startFrame]
	vec := gMap.GetTilePositionAtFeet(e.mTileX, e.mTileY, spriteFrame.W(), spriteFrame.H())
	e.mSprite.Draw(global.gWin, pixel.IM.Moved(vec))
}
