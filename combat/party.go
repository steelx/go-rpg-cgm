package combat

type Party struct {
	Members map[string]*Actor
	world   *WorldExtended
}

func PartyCreate(w *WorldExtended) *Party {
	return &Party{
		Members: make(map[string]*Actor),
		world:   w,
	}
}

func (p *Party) Add(member Actor) {
	p.Members[member.Id] = &member
	p.Members[member.Id].worldRef = p.world
}
func (p *Party) Remove(member Actor) {
	p.removeById(member.Id)
}
func (p *Party) removeById(id string) {
	delete(p.Members, id)
}
