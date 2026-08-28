package adapters

import (
	"context"
	"errors"

	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
)

type WarrantyReader interface {
	Lookup(ctx context.Context, phone, serial string) (*model.WarrantyRecord, error)
}

type WarrantyAdapter struct{ store *persistence.Store }

func NewWarrantyAdapter(store *persistence.Store) *WarrantyAdapter {
	return &WarrantyAdapter{store: store}
}

func (a *WarrantyAdapter) Lookup(ctx context.Context, phone, serial string) (*model.WarrantyRecord, error) {
	if a == nil || a.store == nil {
		return nil, errors.New("warranty adapter is unavailable")
	}
	record, err := a.store.FindWarranty(ctx, phone, serial)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	return record, nil
}

func (a *WarrantyAdapter) LookupActive(ctx context.Context, phone, serial string) (*model.WarrantyRecord, error) {
	record, err := a.Lookup(ctx, phone, serial)
	if err != nil {
		return nil, err
	}
	if record != nil && !record.Active {
		return nil, nil
	}
	return record, nil
}

func (a *WarrantyAdapter) Exists(ctx context.Context, phone, serial string) (bool, error) {
	record, err := a.Lookup(ctx, phone, serial)
	if err != nil {
		return false, err
	}
	return record != nil, nil
}

func (a *WarrantyAdapter) SaveWarranty(ctx context.Context, record model.WarrantyRecord, points []model.ServicePoint) error {
	if a == nil || a.store == nil {
		return errors.New("warranty adapter is unavailable")
	}
	return a.store.SaveWarranty(ctx, record, points)
}

func (a *WarrantyAdapter) ListWarranties(ctx context.Context, activeOnly bool) ([]model.WarrantyRecord, error) {
	if a == nil || a.store == nil {
		return nil, errors.New("warranty adapter is unavailable")
	}
	return a.store.ListWarranties(ctx, activeOnly)
}
