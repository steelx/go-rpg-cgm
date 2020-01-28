package combat

type Party struct {
	Members map[string]*Actor
}

func PartyCreate() *Party {
	return &Party{
		Members: make(map[string]*Actor),
	}
}

func (p *Party) Add(member Actor) {
	p.Members[member.Id] = &member
}
func (p *Party) Remove(member Actor) {
	p.removeById(member.Id)
}
func (p *Party) removeById(id string) {
	delete(p.Members, id)
}
