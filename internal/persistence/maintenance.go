package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"warrantyservice/internal/model"
)

type IntegrityReport struct {
	WarrantyCount int
	PointCount    int
	OrphanLinks   int
	AuditCount    int
	AdviceCount   int
}

func (s *Store) DeactivateWarranty(ctx context.Context, phone, serial, updatedAt string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE warranties SET active=0, updated_at=? WHERE phone=? AND serial_number=?`, updatedAt, phone, serial)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) DeleteWarranty(ctx context.Context, phone, serial string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM warranties WHERE phone=? AND serial_number=?`, phone, serial).Scan(&id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM warranty_points WHERE warranty_id=?`, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM warranties WHERE id=?`, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) Integrity(ctx context.Context) (IntegrityReport, error) {
	var report IntegrityReport
	queries := []struct {
		ptr   *int
		query string
	}{{&report.WarrantyCount, `SELECT COUNT(*) FROM warranties`}, {&report.PointCount, `SELECT COUNT(*) FROM service_points`}, {&report.OrphanLinks, `SELECT COUNT(*) FROM warranty_points w LEFT JOIN warranties r ON r.id=w.warranty_id WHERE r.id IS NULL`}, {&report.AuditCount, `SELECT COUNT(*) FROM admin_audits`}, {&report.AdviceCount, `SELECT COUNT(*) FROM support_advices`}}
	for _, item := range queries {
		if err := s.db.QueryRowContext(ctx, item.query).Scan(item.ptr); err != nil {
			return IntegrityReport{}, err
		}
	}
	return report, nil
}

func (s *Store) ReplacePoints(ctx context.Context, warrantyID int64, points []model.ServicePoint) error {
	if warrantyID <= 0 {
		return errors.New("warranty id must be positive")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM warranty_points WHERE warranty_id=?`, warrantyID); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, point := range points {
		if err := point.Validate(); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := savePointTx(ctx, tx, point); err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO warranty_points(warranty_id,point_id) VALUES(?,?)`, warrantyID, point.ID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) UpsertServicePoint(ctx context.Context, point model.ServicePoint) error {
	if err := point.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO service_points(id,name,address,phone,categories,open_hours) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,address=excluded.address,phone=excluded.phone,categories=excluded.categories,open_hours=excluded.open_hours`, point.ID, point.Name, point.Address, point.Phone, join(point.Categories), point.OpenHours)
	return err
}

func (s *Store) RemoveServicePoint(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("point id must be positive")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM warranty_points WHERE point_id=?`, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM service_points WHERE id=?`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if count == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func (s *Store) Search(ctx context.Context, term string) ([]model.WarrantyRecord, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return s.ListWarranties(ctx, false)
	}
	wildcard := "%" + term + "%"
	rows, err := s.db.QueryContext(ctx, `SELECT id,phone,serial_number,expiry_date,purchase_channel,category,updated_at,active FROM warranties WHERE phone LIKE ? OR serial_number LIKE ? OR category LIKE ? ORDER BY id`, wildcard, wildcard, wildcard)
	if err != nil {
		return nil, fmt.Errorf("search warranties: %w", err)
	}
	defer rows.Close()
	result := []model.WarrantyRecord{}
	for rows.Next() {
		var r model.WarrantyRecord
		var active int
		if err := rows.Scan(&r.ID, &r.Phone, &r.SerialNumber, &r.ExpiryDate, &r.PurchaseChannel, &r.Category, &r.UpdatedAt, &active); err != nil {
			return nil, err
		}
		r.Active = intBool(active)
		r.ServicePoints, err = s.pointsForWarranty(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (s *Store) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	if fn == nil {
		return errors.New("transaction callback is required")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) Vacuum(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `VACUUM`)
	return err
}
