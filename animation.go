package main

type Animation struct {
	mFrames []int
	mIndex  int
	mSPF    float64
	mTime   float64
	mLoop   bool
}

//Animation spf seconds per frame 1.2
func AnimationCreate(frames []int, loop bool, spf float64) Animation {
	return Animation{
		mFrames: frames,
		mIndex:  0,
		mSPF:    spf,
		mTime:   0,
		mLoop:   loop,
	}
}

func (a *Animation) Update(dt float64) {
	// update the animation strip
	a.mTime = a.mTime + dt

	if a.mTime >= a.mSPF {
		a.mIndex += 1
		a.mTime = 0

		if a.IsLastFrame() {
			if a.mLoop {
				a.mIndex = 0
			} else {
				a.mIndex = len(a.mFrames) - 1
			}
		}
	}
}

func (a *Animation) SetFrames(frames []int) {
	a.mFrames = frames
	a.mIndex = 0
}

func (a Animation) Frame() int {
	return a.mFrames[a.mIndex]
}

func (a Animation) GetFirstFrame() int {
	return a.mFrames[0]
}

func (a Animation) IsLastFrame() bool {
	return a.mIndex >= len(a.mFrames)
}

func (a Animation) IsFinished() bool {
	return a.mLoop == false || a.IsLastFrame()
}
