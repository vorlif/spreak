package humanize

import "math"

// LanguageName returns the name of the spoken language as called by the languages used.
func (h *Humanizer) LanguageName(lang string) string {
	return h.loc.Get(lang)
}

func floorDivision(a, b float64) int64 {
	return int64(math.Floor(toFixed(a/b, 3)))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}
