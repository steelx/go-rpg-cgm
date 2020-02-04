package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type MoveState struct {
	Character  *Character
	Map        *GameMap
	Entity     *Entity
	Controller *state_machine.StateMachine
	// ^above common with WaitState
	TileWidth      float64
	MoveX, MoveY   float64
	PixelX, PixelY float64
	MoveSpeed      float64
	Tween          animation.Tween
	Anim           animation.Animation
}

//character *Character, gMap *GameMap
func MoveStateCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	character := charV.Interface().(*Character)
	gMapV := reflect.ValueOf(args[1])
	gMap := gMapV.Interface().(*GameMap)

	s := &MoveState{}
	s.Character = character
	s.Map = gMap
	s.TileWidth = gMap.TileWidth
	s.Entity = character.Entity
	s.Controller = character.Controller
	s.MoveX = 0
	s.MoveY = 0
	s.MoveSpeed = 0.38
	s.Anim = animation.Create([]int{s.Entity.StartFrame}, true, 0.12)
	return s
}

//The StateMachine class requires each state to have
// four functions: Enter, Exit, Render and Update

func (s MoveState) IsFinished() bool {
	return true
}

func (s *MoveState) Enter(data ...interface{}) {
	var frames []int
	v := reflect.ValueOf(data[0])
	dataI := v.Interface().(Direction)
	if dataI.X == -1 {
		frames = s.Character.Anims[CharacterFacingDirection[3]]
		s.Character.SetFacing(3)
	} else if dataI.X == 1 {
		frames = s.Character.Anims[CharacterFacingDirection[1]]
		s.Character.SetFacing(1)
	} else if dataI.Y == -1 {
		frames = s.Character.Anims[CharacterFacingDirection[0]]
		s.Character.SetFacing(0)
	} else if dataI.Y == 1 {
		frames = s.Character.Anims[CharacterFacingDirection[2]]
		s.Character.SetFacing(2)
	}
	s.Anim.SetFrames(frames)

	//save Move X,Y value to used inside Update call
	s.MoveX = dataI.X
	s.MoveY = dataI.Y
	s.PixelX = s.Entity.TileX
	s.PixelY = s.Entity.TileY
	s.Tween = animation.TweenCreate(0, 1, s.MoveSpeed)

	//stop moving if blocking tile
	targetX, targetY := s.Entity.TileX+dataI.X, s.Entity.TileY+dataI.Y

	if player := s.Map.GetEntityAtPos(targetX, targetY); player != nil ||
		s.Map.IsBlockingTile(int(targetX), int(targetY)) {
		s.MoveX = 0
		s.MoveY = 0
		s.Entity.SetFrame(s.Anim.GetFirstFrame())
		s.Controller.Change("wait", Direction{0, 0})
		return
	}
}

func (s MoveState) Exit() {
	//check if an ENTER Trigger exists on given tile coords
	//tileX, tileY := s.Map.GetTileIndex(s.Entity.TileX, s.Entity.TileY)
	tileX, tileY := s.Entity.TileX, s.Entity.TileY
	if trigger := s.Map.GetTrigger(tileX, tileY); trigger.OnEnter != nil {
		trigger.OnEnter(s.Map, s.Entity, tileX, tileY)
	}
}

func (s *MoveState) Render(win *pixelgl.Window) {
	//pending
}

func (s *MoveState) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame() + 1)

	s.Tween.Update(dt)
	value := s.Tween.Value()
	s.Entity.TileX = s.PixelX + value*s.MoveX
	s.Entity.TileY = s.PixelY + value*s.MoveY

	if s.Tween.IsFinished() {
		s.Entity.StartFrame = s.Anim.GetFirstFrame()
		s.Controller.Change(s.Character.DefaultState, Direction{0, 0})
	}
}
