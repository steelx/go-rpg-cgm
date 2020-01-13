package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/globals"
)

var CharacterFacingDirection = [4]string{"up", "right", "down", "left"}

type EntityDefinition struct {
	Texture       pixel.Picture
	Width, Height float64
	StartFrame    int
	TileX, TileY  float64
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
	Children      map[string]*Entity
}

func CreateEntity(def EntityDefinition) *Entity {
	e := &Entity{}

	e.Texture = def.Texture
	e.Frames = globals.LoadAsFrames(def.Texture, def.Width, def.Height)
	e.Sprite = pixel.NewSprite(def.Texture, e.Frames[def.StartFrame])
	e.Width = def.Width
	e.Height = def.Height
	e.TileX = def.TileX
	e.TileY = def.TileY
	e.StartFrame = def.StartFrame
	e.Children = make(map[string]*Entity)
	return e
}

func (e *Entity) AddChild(id string, entity *Entity) {
	e.Children[id] = entity
}
func (e *Entity) RemoveChild(id string) {
	delete(e.Children, id)
}

func (e *Entity) SetTilePos(x, y float64) {
	e.TileX = x
	e.TileY = y
}

func (e *Entity) SetFrame(frame int) {
	e.StartFrame = frame
}

//TeleportAndDraw hero SetTilePos
func (e *Entity) TeleportAndDraw(gMap *GameMap, canvas *pixelgl.Canvas) {
	spriteFrame := e.Frames[e.StartFrame]
	vec := gMap.GetTilePositionAtFeet(e.TileX, e.TileY, spriteFrame.W(), spriteFrame.H())
	e.Sprite = pixel.NewSprite(e.Texture, spriteFrame)
	e.Sprite.Draw(canvas, pixel.IM.Moved(vec))
}

func (e Entity) GetTilePositionOnMap(gMap *GameMap) (vec pixel.Vec) {
	spriteFrame := e.Frames[e.StartFrame]
	vec = gMap.GetTilePositionAtFeet(e.TileX, e.TileY, spriteFrame.W(), spriteFrame.H())
	return
}

//Render will render self + any effects on entity e.g. SleepEntity
func (e *Entity) Render(gMap *GameMap, renderer pixel.Target) {
	//Draw self first
	spriteFrame := e.Frames[e.StartFrame]
	position := e.GetTilePositionOnMap(gMap)
	e.Sprite = pixel.NewSprite(e.Texture, spriteFrame)
	e.Sprite.Draw(renderer, pixel.IM.Moved(position))

	//Draw children
	if len(e.Children) > 0 {
		for _, child := range e.Children {
			spriteFrame := child.Frames[e.StartFrame]
			childPos := pixel.V(child.TileX+position.X, child.TileY+position.Y)
			child.Sprite = pixel.NewSprite(child.Texture, spriteFrame)
			child.Sprite.Draw(renderer, pixel.IM.Moved(childPos))
		}
	}
}

//RenderWithNPC Just had an idea about future renders WIP
//func (e *Entity) RenderWithNPC(renderer pixel.Target) {
//	var others []*Entity
//	for _, npc := range e.NPCs {
//		others = append(others, npc)
//	}
//
//	//sort players as per visible to screen Y position
//	withOthers := append([]*Entity{e}, others...)
//	sort.Slice(withOthers[:], func(i, j int) bool {
//		return withOthers[i].TileY < withOthers[j].TileY
//	})
//
//	for _, player := range withOthers {
//		player.Render(renderer)
//	}
//}
