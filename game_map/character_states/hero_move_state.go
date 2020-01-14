package character_states

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type MoveState struct {
	Character  *game_map.Character
	Map        *game_map.GameMap
	Entity     *game_map.Entity
	Controller *state_machine.StateMachine
	// ^above common with WaitState
	TileWidth      float64
	MoveX, MoveY   float64
	PixelX, PixelY float64
	MoveSpeed      float64
	Tween          animation.Tween
	Anim           animation.Animation
}

func MoveStateCreate(character *game_map.Character, gMap *game_map.GameMap) state_machine.State {
	s := &MoveState{}
	s.Character = character
	s.Map = gMap
	s.TileWidth = gMap.TileWidth
	s.Entity = character.Entity
	s.Controller = character.Controller
	s.MoveX = 0
	s.MoveY = 0
	s.Tween = animation.TweenCreate(0, 0, 1)
	s.MoveSpeed = 0.42
	s.Anim = animation.AnimationCreate([]int{s.Entity.StartFrame}, true, 0.11)
	return s
}

//The StateMachine class requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *MoveState) Enter(dataI interface{}) {
	var frames []int
	v := reflect.ValueOf(dataI)
	data := v.Interface().(globals.Direction)
	if data.X == -1 {
		frames = s.Character.Anims[game_map.CharacterFacingDirection[3]]
		s.Character.SetFacing(3)
	} else if data.X == 1 {
		frames = s.Character.Anims[game_map.CharacterFacingDirection[1]]
		s.Character.SetFacing(1)
	} else if data.Y == -1 {
		frames = s.Character.Anims[game_map.CharacterFacingDirection[0]]
		s.Character.SetFacing(0)
	} else if data.Y == 1 {
		frames = s.Character.Anims[game_map.CharacterFacingDirection[2]]
		s.Character.SetFacing(2)
	}
	s.Anim.SetFrames(frames)

	//save Move X,Y value to used inside Update call
	s.MoveX = data.X
	s.MoveY = data.Y
	s.PixelX = s.Entity.TileX
	s.PixelY = s.Entity.TileY
	s.Tween = animation.TweenCreate(0, 1, s.MoveSpeed)

	//stop moving if blocking tile
	targetX, targetY := s.Entity.TileX+data.X, s.Entity.TileY+data.Y

	if player := s.Map.GetEntityAtPos(targetX, targetY); player != nil ||
		s.Map.IsBlockingTile(int(targetX), int(targetY)) {
		s.MoveX = 0
		s.MoveY = 0
		s.Entity.SetFrame(s.Anim.GetFirstFrame())
		s.Controller.Change("wait", globals.Direction{0, 0})
		return
	}
}

func (s MoveState) Exit() {
	//check if an ENTER Trigger exists on given tile coords
	tileX, tileY := s.Map.GetTileIndex(s.Entity.TileX, s.Entity.TileY)
	if trigger := s.Map.GetTrigger(tileX, tileY); trigger.OnEnter != nil {
		trigger.OnEnter(s.Entity)
	}
}

func (s *MoveState) Render(win *pixelgl.Window) {
	//pending
}

func (s *MoveState) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())

	s.Tween.Update(dt)
	value := s.Tween.Value()
	s.Entity.TileX = s.PixelX + value*s.MoveX
	s.Entity.TileY = s.PixelY + value*s.MoveY

	if s.Tween.IsFinished() {
		s.Controller.Change(s.Character.DefaultState, globals.Direction{0, 0})
	}
}
