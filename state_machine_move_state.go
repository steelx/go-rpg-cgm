package main

type MoveState struct {
	mCharacter  FSMObject
	mMap        GameMap
	mEntity     *Entity
	mController *StateMachine
	// ^above common with WaitState
	mTileWidth       float64
	mMoveX, mMoveY   float64
	mPixelX, mPixelY float64
	mMoveSpeed       float64
	//mTween Tween //TODO tween movement
}

func MoveStateCreate(character FSMObject, gMap GameMap) State {
	s := &MoveState{}
	s.mCharacter = character
	s.mMap = gMap
	s.mTileWidth = gMap.mTileWidth
	s.mEntity = character.mEntity
	s.mController = character.mController
	s.mMoveX = 0
	s.mMoveY = 0
	//s.mTween = TweenCreate(0, 0, 1)
	s.mMoveSpeed = 0.3
	return s
}

//The StateMachine class requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *MoveState) Enter(data Direction) {
	//save Move X,Y value to used inside Update call
	s.mMoveX = data.x
	s.mMoveY = data.y
	s.mPixelX = s.mEntity.mTileX
	s.mPixelY = s.mEntity.mTileY
	//s.mTween = TweenCreate(0, s.mTileWidth, s.mMoveSpeed)
}

func (s *MoveState) Exit() {
	s.mEntity.TeleportAndDraw(s.mMap)
}

func (s *MoveState) Render() {
	//pending
}

func (s *MoveState) Update(dt float64) {
	//s.mTween.Update(dt)
	//value := s.mTween.Value()

	s.mEntity.mTileX = s.mPixelX + s.mMoveX
	s.mEntity.mTileY = s.mPixelY + s.mMoveY

	//if s.mTween.IsFinished() {
	//	s.mController.Change("wait", Direction{0, 0})
	//}
	s.mController.Change("wait", Direction{0, 0})
}
