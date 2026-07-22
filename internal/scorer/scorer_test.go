package scorer

import (
	"finder/internal/model"
	"testing"
	"time"
)

func TestScore(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		content  model.Content
		expected float64
	}{
		{
			name: "video within one week",
			content: model.Content{
				Type:        model.ContentTypeVideo,
				Views:       15000,
				Likes:       1200,
				PublishedAt: now.Add(-3 * 24 * time.Hour),
			},
			expected: round((float64(15000)/1000+float64(1200)/100)*1.5 + 5 + float64(1200)/float64(15000)*10),
		},
		{
			name: "video within one month",
			content: model.Content{
				Type:        model.ContentTypeVideo,
				Views:       25000,
				Likes:       2100,
				PublishedAt: now.Add(-15 * 24 * time.Hour),
			},
			expected: round((float64(25000)/1000+float64(2100)/100)*1.5 + 3 + float64(2100)/float64(25000)*10),
		},
		{
			name: "video within three months",
			content: model.Content{
				Type:        model.ContentTypeVideo,
				Views:       50000,
				Likes:       3000,
				PublishedAt: now.Add(-60 * 24 * time.Hour),
			},
			expected: round((float64(50000)/1000+float64(3000)/100)*1.5 + 1 + float64(3000)/float64(50000)*10),
		},
		{
			name: "video older than three months",
			content: model.Content{
				Type:        model.ContentTypeVideo,
				Views:       100000,
				Likes:       5000,
				PublishedAt: now.Add(-120 * 24 * time.Hour),
			},
			expected: round((float64(100000)/1000+float64(5000)/100)*1.5 + 0 + float64(5000)/float64(100000)*10),
		},
		{
			name: "text within one week",
			content: model.Content{
				Type:        model.ContentTypeText,
				ReadingTime: 8,
				Reactions:   320,
				PublishedAt: now.Add(-5 * 24 * time.Hour),
			},
			expected: round((float64(8)+float64(320)/50)*1.0 + 5 + float64(320)/float64(8)*5),
		},
		{
			name: "text within one month",
			content: model.Content{
				Type:        model.ContentTypeText,
				ReadingTime: 12,
				Reactions:   540,
				PublishedAt: now.Add(-20 * 24 * time.Hour),
			},
			expected: round((float64(12)+float64(540)/50)*1.0 + 3 + float64(540)/float64(12)*5),
		},
		{
			name: "video with zero views",
			content: model.Content{
				Type:        model.ContentTypeVideo,
				Views:       0,
				Likes:       0,
				PublishedAt: now.Add(-2 * 24 * time.Hour),
			},
			expected: round(0*1.5 + 5 + 0),
		},
		{
			name: "text with zero reading time",
			content: model.Content{
				Type:        model.ContentTypeText,
				ReadingTime: 0,
				Reactions:   100,
				PublishedAt: now.Add(-10 * 24 * time.Hour),
			},
			expected: round((float64(0)+float64(100)/50)*1.0 + 3 + 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Score(tt.content)
			if got != tt.expected {
				t.Errorf("got %.2f, want %.2f", got, tt.expected)
			}
		})
	}
}
