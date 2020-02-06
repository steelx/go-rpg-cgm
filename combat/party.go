package combat

type Party struct {
	Members map[string]*Actor
	World   *WorldExtended
}

func PartyCreate(w *WorldExtended) *Party {
	return &Party{
		Members: make(map[string]*Actor),
		World:   w,
	}
}

func (p *Party) Add(member Actor) {
	p.Members[member.Id] = &member
	p.Members[member.Id].worldRef = p.World
}
func (p *Party) Remove(member Actor) {
	p.removeById(member.Id)
}
func (p *Party) removeById(id string) {
	delete(p.Members, id)
}

func (p Party) ToArray() []*Actor {
	var party []*Actor
	for _, v := range p.Members {
		party = append(party, v)
	}
	return party
}
