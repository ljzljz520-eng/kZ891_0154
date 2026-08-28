package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

type Health struct {
	Ready  bool   `json:"ready"`
	Driver string `json:"driver"`
	Path   string `json:"path"`
	Tables int    `json:"tables"`
}

func (s *Store) Health(ctx context.Context) Health {
	result := Health{Driver: "modernc.org/sqlite"}
	if s == nil || s.db == nil {
		return result
	}
	result.Path = s.path
	if err := s.db.PingContext(ctx); err != nil {
		return result
	}
	var tables int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tables); err != nil {
		return result
	}
	result.Tables = tables
	result.Ready = tables >= 5
	return result
}

func (s *Store) EnsureReady(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}
	health := s.Health(ctx)
	if !health.Ready {
		return fmt.Errorf("store schema is incomplete")
	}
	return nil
}

func (s *Store) LastInsertID(ctx context.Context) (int64, error) {
	var id sql.NullInt64
	err := s.db.QueryRowContext(ctx, `SELECT MAX(id) FROM admin_audits`).Scan(&id)
	if err != nil {
		return 0, err
	}
	if !id.Valid {
		return 0, nil
	}
	return id.Int64, nil
}
