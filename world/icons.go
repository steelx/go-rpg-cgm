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
	}

	ico.UVs = utilz.LoadAsFramesFromTop(ico.Texture, 16, 16)
	ico.Sprites = make([]*pixel.Sprite, len(ico.UVs))

	for k := range ico.UVs {
		sprite := pixel.NewSprite(ico.Texture, ico.UVs[k])
		ico.Sprites[k] = sprite
	}

	return ico
}

//Get accepts ItemType int e.g. weapon = 3
func (i Icons) Get(d ItemType) *pixel.Sprite {
	return i.Sprites[d]
}
