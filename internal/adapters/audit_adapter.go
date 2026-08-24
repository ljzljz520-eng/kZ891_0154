package adapters

import (
	"context"
	"errors"

	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
)

type AuditWriter struct{ store *persistence.Store }

func NewAuditWriter(store *persistence.Store) *AuditWriter { return &AuditWriter{store: store} }

func (w *AuditWriter) Record(ctx context.Context, audit model.AdminAudit) error {
	if w == nil || w.store == nil {
		return errors.New("audit writer is unavailable")
	}
	return w.store.SaveAudit(ctx, audit)
}

func (w *AuditWriter) History(ctx context.Context, recordID int64) ([]model.AdminAudit, error) {
	if w == nil || w.store == nil {
		return nil, errors.New("audit writer is unavailable")
	}
	return w.store.ListAudits(ctx, recordID)
}
