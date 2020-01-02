package main

type Trigger struct {
	OnEnter func(entity *Entity)
	OnExit  func()
	OnUse   func()
}

//TriggerCreate
//e.g. CastleRoomMap.mTriggers[[2]float64{tileX, tileY}].OnEnter(gHero.mEntity)
func TriggerCreate(OnEnter func(entity *Entity), OnExit func(), OnUse func()) Trigger {
	return Trigger{
		OnEnter: OnEnter,
		OnExit:  OnExit,
		OnUse:   OnUse,
	}
}
