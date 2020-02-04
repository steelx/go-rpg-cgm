package animation

type Animation struct {
	Frames []int
	Index  int
	spf    float64
	time   float64
	loop   bool
}

//Animation spf seconds per frame 1.2
func Create(frames []int, loop bool, spf float64) Animation {
	return Animation{
		Frames: frames,
		Index:  0,
		spf:    spf,
		time:   0,
		loop:   loop,
	}
}

func (a *Animation) Update(dt float64) {
	// update the animation strip
	a.time = a.time + dt

	if a.time >= a.spf {
		a.Index += 1
		a.time = 0

		if a.IsLastFrame() {
			if a.loop {
				a.Index = 0
			} else {
				a.Index = len(a.Frames) - 1
			}
		}
	}
}

func (a *Animation) SetFrames(frames []int) {
	a.Frames = frames
	a.Index = 0
}

func (a Animation) Frame() int {
	return a.Frames[a.Index]
}

func (a Animation) GetFirstFrame() int {
	return a.Frames[0]
}

func (a Animation) IsLastFrame() bool {
	return a.Index >= len(a.Frames)-1
}

func (a Animation) IsFinished() bool {
	return a.loop == false && a.IsLastFrame()
}
