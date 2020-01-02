package main

type UV struct {
	ux, uy, vx, vy float64
}

type State interface {
	Enter(data Direction)
	Render()
	Exit()
	Update(dt float64)
}

type CharacterStateBase struct {
	mCharacter  *Character
	mMap        *GameMap
	mEntity     *Entity
	mController *StateMachine
}

type Direction struct {
	x, y float64
}
