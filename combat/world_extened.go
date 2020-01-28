package combat

import "github.com/steelx/go-rpg-cgm/world"

type WorldExtended struct {
	world.World
	Party *Party
}

func WorldExtendedCreate() *WorldExtended {
	w := &WorldExtended{}
	w.Party = PartyCreate(w)
	w.Time = 0
	w.Gold = 0
	w.Items = make([]world.ItemIndex, 0)
	w.KeyItems = make([]world.ItemIndex, 0)
	w.Icons = world.IconsCreate()
	return w
}
