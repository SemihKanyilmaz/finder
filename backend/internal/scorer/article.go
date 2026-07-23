package scorer

import "finder/internal/model"

type articleScorer struct{}

func (articleScorer) Base(c model.Content) float64 {
	return float64(c.ReadingTime) + float64(c.Reactions)/50
}

func (articleScorer) Coefficient() float64 {
	return 1.0
}

func (articleScorer) Engagement(c model.Content) float64 {
	if c.ReadingTime == 0 {
		return 0
	}
	return float64(c.Reactions) / float64(c.ReadingTime) * 5
}
