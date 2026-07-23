package scorer

import (
	"finder/internal/model"
	"math"
	"time"
)

type calculator interface {
	Base(c model.Content) float64
	Coefficient() float64
	Engagement(c model.Content) float64
}

var calculators = map[model.ContentType]calculator{
	model.ContentTypeVideo:   videoScorer{},
	model.ContentTypeArticle: articleScorer{},
}

func Score(c model.Content) float64 {
	s, ok := calculators[c.Type]
	if !ok {
		return 0
	}
	base := s.Base(c)
	coeff := s.Coefficient()
	freshness := freshnessScore(c.PublishedAt)
	engagement := s.Engagement(c)
	return round(base*coeff + freshness + engagement)
}

func freshnessScore(publishedAt time.Time) float64 {
	age := time.Since(publishedAt)
	switch {
	case age <= 7*24*time.Hour:
		return 5
	case age <= 30*24*time.Hour:
		return 3
	case age <= 90*24*time.Hour:
		return 1
	default:
		return 0
	}
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}
