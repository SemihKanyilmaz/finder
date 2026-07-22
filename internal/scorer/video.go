package scorer

import "finder/internal/model"

type videoScorer struct{}

func (videoScorer) Base(c model.Content) float64 {
	return float64(c.Views)/1000 + float64(c.Likes)/100
}

func (videoScorer) Coefficient() float64 {
	return 1.5
}

func (videoScorer) Engagement(c model.Content) float64 {
	if c.Views == 0 {
		return 0
	}
	return float64(c.Likes) / float64(c.Views) * 10
}
