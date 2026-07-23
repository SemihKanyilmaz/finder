package repository

import (
	"context"
	"finder/internal/model"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentRepository interface {
	Upsert(ctx context.Context, contents []model.Content) error
	Search(ctx context.Context, params model.SearchParams) (model.SearchResult, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) ContentRepository {
	return &postgresRepository{pool: pool}
}

// buildPrefixQuery converts a user keyword into a tsquery that supports prefix matching.
// e.g. "clean arch" → "clean:* & arch:*"
func buildPrefixQuery(keyword string) string {
	words := strings.Fields(strings.ToLower(keyword))
	for i, w := range words {
		words[i] = w + ":*"
	}
	return strings.Join(words, " & ")
}

type queryBuilder struct {
	args []any
}

func (b *queryBuilder) Add(val any) string {
	b.args = append(b.args, val)
	return fmt.Sprintf("$%d", len(b.args))
}

func (r *postgresRepository) Upsert(ctx context.Context, contents []model.Content) error {
	if len(contents) == 0 {
		return nil
	}

	var (
		placeholders []string
		args         []any
	)

	for i, c := range contents {
		base := i * 10
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9, base+10,
		))
		args = append(args, c.ID, c.Source, c.Title, string(c.Type), c.PublishedAt, c.Views, c.Likes, c.ReadingTime, c.Reactions, c.Score)
	}

	query := fmt.Sprintf(`
		INSERT INTO contents (id, source, title, type, published_at, views, likes, reading_time, reactions, score)
		VALUES %s
		ON CONFLICT (id, source) DO UPDATE SET
			title = EXCLUDED.title,
			type = EXCLUDED.type,
			published_at = EXCLUDED.published_at,
			views = EXCLUDED.views,
			likes = EXCLUDED.likes,
			reading_time = EXCLUDED.reading_time,
			reactions = EXCLUDED.reactions,
			score = EXCLUDED.score,
			updated_at = NOW()`,
		strings.Join(placeholders, ", "),
	)

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("upsert: %w", err)
	}

	return nil
}

func (r *postgresRepository) Search(ctx context.Context, params model.SearchParams) (model.SearchResult, error) {
	var conditions []string
	qb := &queryBuilder{}

	if params.Keyword != "" {
		conditions = append(conditions, fmt.Sprintf(
			"search_vector @@ to_tsquery('simple', %s)", qb.Add(buildPrefixQuery(params.Keyword)),
		))
	}

	if params.ContentType != "" {
		conditions = append(conditions, fmt.Sprintf("type = %s", qb.Add(string(params.ContentType))))
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "ORDER BY score DESC"
	if params.SortBy == "freshness" {
		orderBy = "ORDER BY published_at DESC"
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	query := fmt.Sprintf(
		`SELECT id, source, title, type, published_at, views, likes, reading_time, reactions, score, COUNT(*) OVER() AS total_count
		FROM contents %s %s LIMIT %s OFFSET %s`,
		where, orderBy, qb.Add(params.PageSize), qb.Add(offset),
	)

	rows, err := r.pool.Query(ctx, query, qb.args...)
	if err != nil {
		return model.SearchResult{}, fmt.Errorf("search query: %w", err)
	}
	defer rows.Close()

	var (
		items      = make([]model.Content, 0, params.PageSize)
		totalCount int
	)

	for rows.Next() {
		var c model.Content
		var contentType string
		if err := rows.Scan(&c.ID, &c.Source, &c.Title, &contentType, &c.PublishedAt, &c.Views, &c.Likes, &c.ReadingTime, &c.Reactions, &c.Score, &totalCount); err != nil {
			return model.SearchResult{}, fmt.Errorf("scan row: %w", err)
		}
		c.Type = model.ContentType(contentType)
		items = append(items, c)
	}

	if err := rows.Err(); err != nil {
		return model.SearchResult{}, fmt.Errorf("rows iteration: %w", err)
	}

	totalPages := 0
	if params.PageSize > 0 {
		totalPages = (totalCount + params.PageSize - 1) / params.PageSize
	}

	return model.SearchResult{
		Items:      items,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}
