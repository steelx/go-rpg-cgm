package game_states

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/state_stacks"
	"log"
)

type ExploreState struct {
	Stack  state_stacks.StateStack
	MapDef *tilepix.Map
	Map    *game_map.GameMap
	Hero   *game_map.Character
}

func ExploreStateCreate(stack state_stacks.StateStack, tilemap *tilepix.Map, startPos pixel.Vec) ExploreState {
	es := ExploreState{
		Stack:  stack,
		MapDef: tilemap,
	}

	es.Map = game_map.MapCreate(es.MapDef)

	heroPng, err := globals.LoadPicture("../resources/walk_cycle.png")
	if err != nil {
		log.Fatal(err)
	}

	es.Hero = game_map.CharacterCreate("Ajinkya",
		[][]int{{16, 17, 18, 19}, {20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30, 31}},
		game_map.CharacterFacingDirection[2],
		game_map.CharacterDefinition{
			Texture: heroPng, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      startPos.X,
			TileY:      startPos.Y,
			Map:        es.Map,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return character_states.WaitStateCreate(es.Hero, es.Map)
			},
			"move": func() state_machine.State {
				return character_states.MoveStateCreate(es.Hero, es.Map)
			},
		},
	)
	es.Map.GoToTile(startPos.X, startPos.Y)

	return es
}

func (es ExploreState) Enter(data globals.Direction) {}
func (es ExploreState) Exit()                        {}
func (es ExploreState) Update(dt float64)            {}
func (es ExploreState) Render()                      {}
func (es ExploreState) HandleInput()                 {}
