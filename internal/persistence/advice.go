package persistence

import (
	"context"
	"database/sql"
	"errors"
	"warrantyservice/internal/model"
)

func (s *Store) SaveAdvice(ctx context.Context, advice model.SupportAdvice) error {
	if advice.Phone == "" || advice.SerialNumber == "" || advice.Message == "" {
		return errors.New("advice fields are required")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO support_advices(phone,serial_number,message,created_at,hotline,channels) VALUES(?,?,?,?,?,?)`, advice.Phone, advice.SerialNumber, advice.Message, advice.CreatedAt, advice.Hotline, advice.Channels)
	return err
}

func (s *Store) LatestAdvice(ctx context.Context, phone, serial string) (*model.SupportAdvice, error) {
	var a model.SupportAdvice
	err := s.db.QueryRowContext(ctx, `SELECT id,phone,serial_number,message,created_at,hotline,channels FROM support_advices WHERE phone=? AND serial_number=? ORDER BY id DESC LIMIT 1`, phone, serial).Scan(&a.ID, &a.Phone, &a.SerialNumber, &a.Message, &a.CreatedAt, &a.Hotline, &a.Channels)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) Count(ctx context.Context, table string) (int, error) {
	allowed := map[string]bool{"warranties": true, "admin_audits": true, "support_advices": true, "service_points": true}
	if !allowed[table] {
		return 0, errors.New("unsupported table")
	}
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
