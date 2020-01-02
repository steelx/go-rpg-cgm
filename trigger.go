package main

type Trigger struct {
	OnEnter func(entity *Entity)
	OnExit  func()
	OnUse   func()
}

//TriggerCreate
//e.g. tileX, tileY = CastleRoomMap.GetTileIndex(9, 10)
//     CastleRoomMap.mTriggers[[2]float64{tileX, tileY}].OnEnter(gHero.mEntity)
func TriggerCreate(OnEnter func(entity *Entity), OnExit func(), OnUse func()) Trigger {
	//OnUse: When the spacebar is pressed,
	// -> the tile that the character is facing is checked for triggers
	//OnEnter: When user walkover a tile -> trigger is executed
	//OnExit: when user stop moving on a tile -> trigger is executed
	return Trigger{
		OnEnter: OnEnter,
		OnExit:  OnExit,
		OnUse:   OnUse,
	}
}
