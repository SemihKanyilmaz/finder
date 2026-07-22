package scorer

import "finder/internal/model"

type textScorer struct{}

func (textScorer) Base(c model.Content) float64 {
	return float64(c.ReadingTime) + float64(c.Reactions)/50
}

func (textScorer) Coefficient() float64 {
	return 1.0
}

func (textScorer) Engagement(c model.Content) float64 {
	if c.ReadingTime == 0 {
		return 0
	}
	return float64(c.Reactions) / float64(c.ReadingTime) * 5
}
