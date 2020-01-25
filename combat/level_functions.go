package combat

import "math"

//exponent represents the increase of difficulty between levels
//Disgea Level formula to the first 99 levels.
func NextLevel(level int) float64 {
	exponent := 1
	baseXP := 1000.0
	//return math.Round( 0.04 * float64(level ^ 3) + 0.8 * float64(level ^ 2) + float64(2 * level))
	return math.Round(baseXP*float64(level^exponent) + 0.8*float64(level^2) + float64(2*level))
}

//func NextLevel_(level int) float64 {
//	exponent := 2
//	baseXP := 1000.0
//	return math.Floor(baseXP * float64(level^exponent))
//}
