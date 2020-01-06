package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/globals"
)

var CharacterFacingDirection = [4]string{"up", "right", "down", "left"}

type CharacterDefinition struct {
	Texture       pixel.Picture
	Width, Height float64
	StartFrame    int
	TileX, TileY  float64
	Map           *GameMap
}

//Entity represents any kind of map object from a
// treasure chest to an NPC
type Entity struct {
	Sprite        *pixel.Sprite
	Texture       pixel.Picture
	Height, Width float64
	TileX, TileY  float64
	StartFrame    int
	Frames        []pixel.Rect
	Map           *GameMap
}

func CreateEntity(def CharacterDefinition) *Entity {
	e := &Entity{}

	e.Map = def.Map
	e.Texture = def.Texture
	e.Frames = globals.LoadAsFrames(def.Texture, def.Width, def.Height)
	e.Sprite = pixel.NewSprite(def.Texture, e.Frames[def.StartFrame])
	e.Width = def.Width
	e.Height = def.Height
	e.TileX = def.TileX
	e.TileY = def.TileY
	e.StartFrame = def.StartFrame
	return e
}

func (e *Entity) SetFrame(frame int) {
	e.StartFrame = frame
}

//TeleportAndDraw hero movement & set position for sprite
func (e *Entity) TeleportAndDraw(gMap GameMap, canvas *pixelgl.Canvas) {
	spriteFrame := e.Frames[e.StartFrame]
	vec := gMap.GetTilePositionAtFeet(e.TileX, e.TileY, spriteFrame.W(), spriteFrame.H())
	e.Sprite = pixel.NewSprite(e.Texture, spriteFrame)
	e.Sprite.Draw(canvas, pixel.IM.Moved(vec))
}
