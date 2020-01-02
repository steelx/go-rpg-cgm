package main

var gHero *Character

func init() {
	pic, err := LoadPicture("./resources/walk_cycle.png")
	panicIfErr(err)

	gHero = &Character{
		mAnimUp:    []int{16, 17, 18, 19},
		mAnimRight: []int{20, 21, 22, 23},
		mAnimDown:  []int{24, 25, 26, 27},
		mAnimLeft:  []int{28, 29, 30, 31},
		mFacing:    CharacterFacingDirection[2],
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 24,
			tileX:      4,
			tileY:      4,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return WaitStateCreate(gHero, CastleRoomMap)
				},
				"move": func() State {
					return MoveStateCreate(gHero, CastleRoomMap)
				},
			},
		),
	}
}
