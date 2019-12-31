package main

type Tween struct {
	isFinished                                               bool
	tweenF                                                   func(timePassed, start, distance, duration float64) float64
	distance, startValue, current, totalDuration, timePassed float64
}

//@start				start value
//@finish				end value
//@totalDuration  		time in which to perform tween.
//@tweenF				tween function defaults to linear
func TweenCreate(start, finish, totalDuration float64) Tween {
	return Tween{
		tweenF:        TweenLinear,
		distance:      finish - start,
		startValue:    start,
		current:       start,
		totalDuration: totalDuration,
		timePassed:    0,
		isFinished:    false,
	}
}

func (tw Tween) IsFinished() bool {
	return tw.isFinished
}

func (tw Tween) Value() float64 {
	return tw.current
}

func TweenLinear(timePassed, start, distance, duration float64) float64 {
	return distance*timePassed/duration + start
}

func (tw Tween) FinishValue() float64 {
	return tw.startValue + tw.distance
}

func (tw *Tween) Update(elapsedTime float64) {
	tw.timePassed = tw.timePassed + elapsedTime //delta time
	tw.current = tw.tweenF(tw.timePassed, tw.startValue, tw.distance, tw.totalDuration)

	if tw.timePassed > tw.totalDuration {
		tw.current = tw.startValue + tw.distance
		tw.isFinished = true
	}
}
