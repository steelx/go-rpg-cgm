package world

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/utilz"
)

type IconDefs struct {
	Usable,
	Accessory,
	Weapon,
	Sword,
	Dagger,
	Stave,
	Armor,
	Plate,
	Leather,
	Robe,
	UpArrow,
	DownArrow int
}

type Icons struct {
	Texture  pixel.Picture
	UVs      []pixel.Rect
	Sprites  []*pixel.Sprite
	IconDefs IconDefs
}

func IconsCreate() Icons {
	inventoryIconsPng, err := utilz.LoadPicture("../resources/inventory_icons.png")
	utilz.PanicIfErr(err)
	//488
	ico := Icons{
		Texture: inventoryIconsPng,
		IconDefs: IconDefs{
			Usable:    1,
			Accessory: 2,
			Weapon:    3,
			Sword:     4,
			Dagger:    5,
			Stave:     6,
			Armor:     7,
			Plate:     8,
			Leather:   9,
			Robe:      10,
			UpArrow:   11,
			DownArrow: 12,
		},
	}

	ico.UVs = utilz.LoadAsFramesFromTop(ico.Texture, 18, 18)
	ico.Sprites = make([]*pixel.Sprite, len(ico.UVs))

	for k := range ico.UVs {
		sprite := pixel.NewSprite(ico.Texture, ico.UVs[k])
		ico.Sprites[k] = sprite
	}

	return ico
}

//Get accepts ItemType int e.g. weapon = 3
func (i Icons) Get(d int) *pixel.Sprite {
	return i.Sprites[d-1]
}
