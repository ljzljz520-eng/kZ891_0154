package adapters

import (
	"context"
	"errors"
	"sort"
	"strings"

	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
)

type Catalog struct{ store *persistence.Store }

func NewCatalog(store *persistence.Store) *Catalog { return &Catalog{store: store} }

func (c *Catalog) All(ctx context.Context) ([]model.ServicePoint, error) {
	if c == nil || c.store == nil {
		return nil, errors.New("catalog is unavailable")
	}
	count, err := c.store.Count(ctx, "service_points")
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return []model.ServicePoint{}, nil
	}
	return c.store.FindServicePoints(ctx, c.names(ctx))
}

func (c *Catalog) names(ctx context.Context) []string {
	result := []string{}
	if c == nil || c.store == nil {
		return result
	}
	rows, err := c.store.DB().QueryContext(ctx, `SELECT name FROM service_points ORDER BY name`)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			result = append(result, name)
		}
	}
	return result
}

func (c *Catalog) FindByCategory(ctx context.Context, category string) ([]model.ServicePoint, error) {
	points, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	category = strings.TrimSpace(category)
	if category == "" {
		return points, nil
	}
	result := make([]model.ServicePoint, 0)
	for _, point := range points {
		for _, candidate := range point.Categories {
			if candidate == category {
				result = append(result, point)
				break
			}
		}
	}
	return result, nil
}

func (c *Catalog) Categories(ctx context.Context) ([]string, error) {
	points, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for _, point := range points {
		for _, category := range point.Categories {
			if category != "" {
				set[category] = true
			}
		}
	}
	result := make([]string, 0, len(set))
	for category := range set {
		result = append(result, category)
	}
	sort.Strings(result)
	return result, nil
}

func (c *Catalog) HasPoint(ctx context.Context, name string) (bool, error) {
	if c == nil || c.store == nil {
		return false, errors.New("catalog is unavailable")
	}
	rows, err := c.store.DB().QueryContext(ctx, `SELECT 1 FROM service_points WHERE name=? LIMIT 1`, name)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
