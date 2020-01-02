package main

type Character struct {
	mAnimUp     []int
	mAnimRight  []int
	mAnimDown   []int
	mAnimLeft   []int
	mFacing     string
	mEntity     *Entity
	mController *StateMachine //[name] -> [function that returns state]
}

func (ch Character) GetFacedTileCoords() (x, y float64) {
	var xOff, yOff float64 = 0, 0
	if ch.mFacing == CharacterFacingDirection[3] {
		xOff = -1 //"left"
	} else if ch.mFacing == CharacterFacingDirection[1] {
		xOff = 1 //"right"
	} else if ch.mFacing == CharacterFacingDirection[0] {
		yOff = -1 //"up"
	} else if ch.mFacing == CharacterFacingDirection[2] {
		yOff = 1 //"down"
	}

	x = ch.mEntity.mTileX + xOff
	y = ch.mEntity.mTileY + yOff
	return
}

func (ch *Character) SetFacing(dir int) {
	ch.mFacing = CharacterFacingDirection[dir]
}
