package game_map

import (
	"github.com/faiface/pixel/pixelgl"
)

type BlockUntilEvent struct {
	UntilFunc func() bool
}

func BlockUntilEventCreate(untilFunc func() bool) *BlockUntilEvent {

	return &BlockUntilEvent{
		UntilFunc: untilFunc,
	}
}

func (b BlockUntilEvent) Update(dt float64) {
}

func (b BlockUntilEvent) IsBlocking() bool {
	return !b.UntilFunc()
}

func (b BlockUntilEvent) IsFinished() bool {
	return !b.IsBlocking()
}

func (b BlockUntilEvent) Render(win *pixelgl.Window) {
}
