package world

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/utilz"
)

//var IconPNGs Icons
//
//func init()  {
//	inventoryIconsPng, err := globals.LoadPicture("../resources/inventory_icons.png")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	IconPNGs = IconsCreate(inventoryIconsPng)
//}

type IconDefs int

const (
	UsableICO IconDefs = iota
	AccessoryICO
	WeaponICO
	ArmorICO
	UpArrowICO
	OwnArrowICO
)

func (d IconDefs) String() string {
	return [...]string{"UsableICO", "AccessoryICO", "WeaponICO", "ArmorICO", "UpArrowICO", "DownArrowICO"}[d]
}

type Icons struct {
	Texture pixel.Picture
	UVs     []pixel.Rect
	Sprites []*pixel.Sprite
}

func IconsCreate(pic pixel.Picture) Icons {
	i := Icons{
		Texture: pic,
	}

	i.UVs = utilz.LoadAsFrames(i.Texture, 18, 18)

	for k := range i.UVs {
		sprite := pixel.NewSprite(i.Texture, i.UVs[k])
		i.Sprites[k] = sprite
	}

	return i
}

func (i Icons) Get(d ItemType) *pixel.Sprite {
	return i.Sprites[d]
}
