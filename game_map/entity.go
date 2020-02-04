package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/utilz"
)

var CharacterFacingDirection = [4]string{"up", "right", "down", "left"}

type EntityDefinition struct {
	Texture       string
	Width, Height float64
	StartFrame    int
	Frames        []int //used for quick FX
	TileX, TileY  float64
}

//Entity represents any kind of map object from a
// treasure chest to an NPC
type Entity struct {
	Sprite        *pixel.Sprite
	Texture       pixel.Picture
	Height, Width float64
	TileX, TileY  float64
	X, Y          float64 //used with Combat background image
	StartFrame    int
	Frames        []pixel.Rect
	Children      map[string]*Entity
}

func CreateEntity(def EntityDefinition) *Entity {
	pngImg, err := utilz.LoadPicture(def.Texture)
	utilz.PanicIfErr(err)
	e := &Entity{}

	e.Texture = pngImg
	e.Frames = utilz.LoadAsFrames(pngImg, def.Width, def.Height)
	e.Sprite = pixel.NewSprite(pngImg, e.Frames[def.StartFrame])
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
func (e *Entity) Render(gMap *GameMap, renderer pixel.Target, pos pixel.Vec) {
	//Draw self first
	spriteFrame := e.Frames[e.StartFrame]
	position := pixel.ZV
	if gMap != nil {
		position = e.GetTilePositionOnMap(gMap)
	} else {
		position = pos
	}
	e.Sprite = pixel.NewSprite(e.Texture, spriteFrame)
	e.Sprite.Draw(renderer, pixel.IM.Moved(position))
	//Draw children
	if len(e.Children) > 0 {
		for _, child := range e.Children {
			spriteFrame := child.Frames[child.StartFrame]
			childPos := pixel.V(child.TileX+position.X, child.TileY+position.Y)
			child.Sprite = pixel.NewSprite(child.Texture, spriteFrame)
			child.Sprite.Draw(renderer, pixel.IM.Moved(childPos))
		}
	}
}

//GetSelectPosition gets Head position minus offset
func (e *Entity) GetSelectPosition() pixel.Vec {
	x := e.X
	y := e.Y + (e.Height / 2) + 10
	return pixel.V(x, y)
}

func (e *Entity) GetTargetPosition() pixel.Vec {
	x := e.X - (e.Width / 2) - 20
	y := e.Y - (e.Height / 2)
	return pixel.V(x, y)
}
