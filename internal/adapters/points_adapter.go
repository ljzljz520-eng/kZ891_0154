package adapters

import (
	"context"
	"errors"

	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
)

type PointDirectory struct{ store *persistence.Store }

func NewPointDirectory(store *persistence.Store) *PointDirectory {
	return &PointDirectory{store: store}
}

func (d *PointDirectory) Resolve(ctx context.Context, names []string) ([]model.ServicePoint, error) {
	if d == nil || d.store == nil {
		return nil, errors.New("point directory is unavailable")
	}
	return d.store.FindServicePoints(ctx, names)
}

func (d *PointDirectory) FilterByCategory(points []model.ServicePoint, category string) []model.ServicePoint {
	if category == "" {
		return points
	}
	result := make([]model.ServicePoint, 0, len(points))
	for _, point := range points {
		for _, candidate := range point.Categories {
			if candidate == category {
				result = append(result, point)
				break
			}
		}
	}
	return result
}

func (d *PointDirectory) Summaries(points []model.ServicePoint) []string {
	result := make([]string, 0, len(points))
	for _, point := range points {
		result = append(result, point.DisplayName()+"（"+point.Phone+"）")
	}
	return result
}
