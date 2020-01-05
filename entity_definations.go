package main

var (
	gHero *Character
	gNPC1 *Character
	gNPC2 *Character
)

func init() {
	pic, err := LoadPicture("./resources/walk_cycle.png")
	panicIfErr(err)

	gHero = &Character{
		name:       "Ajinkya",
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
			gMap:       CastleRoomMap,
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

	gNPC1 = &Character{
		name:    "Aghori Baba",
		mFacing: CharacterFacingDirection[2],
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 46,
			tileX:      9,
			tileY:      4,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return NPCWaitStateCreate(gNPC1, CastleRoomMap)
				},
			},
		),
	}

	gNPC2 = &Character{
		name:       "Bhadrasaal",
		mAnimUp:    []int{48, 49, 50, 51},
		mAnimRight: []int{52, 53, 54, 55},
		mAnimDown:  []int{56, 57, 58, 59},
		mAnimLeft:  []int{60, 61, 62, 63},
		mFacing:    CharacterFacingDirection[2],
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 56,
			tileX:      3,
			tileY:      8,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return NPCStrollWaitStateCreate(gNPC2, CastleRoomMap)
				},
				"move": func() State {
					return MoveStateCreate(gNPC2, CastleRoomMap)
				},
			},
		),
	}

	//Init Characters
	gHero.mController.Change("wait", Direction{0, 0})
	gNPC1.mController.Change("wait", Direction{0, 0})
	gNPC2.mController.Change("wait", Direction{0, 0})
}
