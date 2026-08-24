package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"warrantyservice/internal/model"
)

func (s *Store) FindWarranty(ctx context.Context, phone, serial string) (*model.WarrantyRecord, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id,phone,serial_number,expiry_date,purchase_channel,category,updated_at,active FROM warranties WHERE phone=? AND serial_number=?`, phone, serial)
	var record model.WarrantyRecord
	var active int
	if err := row.Scan(&record.ID, &record.Phone, &record.SerialNumber, &record.ExpiryDate, &record.PurchaseChannel, &record.Category, &record.UpdatedAt, &active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find warranty: %w", err)
	}
	record.Active = intBool(active)
	points, err := s.pointsForWarranty(ctx, record.ID)
	if err != nil {
		return nil, err
	}
	record.ServicePoints = points
	return &record, nil
}

func (s *Store) pointsForWarranty(ctx context.Context, id int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT p.name FROM service_points p JOIN warranty_points w ON p.id=w.point_id WHERE w.warranty_id=? ORDER BY p.name`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) FindServicePoints(ctx context.Context, names []string) ([]model.ServicePoint, error) {
	if len(names) == 0 {
		return []model.ServicePoint{}, nil
	}
	result := make([]model.ServicePoint, 0, len(names))
	for _, name := range names {
		var point model.ServicePoint
		var categories string
		err := s.db.QueryRowContext(ctx, `SELECT id,name,address,phone,categories,open_hours FROM service_points WHERE name=?`, name).Scan(&point.ID, &point.Name, &point.Address, &point.Phone, &categories, &point.OpenHours)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, err
		}
		point.Categories = split(categories)
		result = append(result, point)
	}
	return result, nil
}

func (s *Store) ListWarranties(ctx context.Context, activeOnly bool) ([]model.WarrantyRecord, error) {
	query := `SELECT id,phone,serial_number,expiry_date,purchase_channel,category,updated_at,active FROM warranties`
	if activeOnly {
		query += ` WHERE active=1`
	}
	query += ` ORDER BY id`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.WarrantyRecord
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

func (s *Store) SaveAudit(ctx context.Context, audit model.AdminAudit) error {
	if audit.RecordID == 0 || audit.Action == "" || audit.Operator == "" {
		return errors.New("audit fields are required")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO admin_audits(record_id,action,operator,occurred_at,detail) VALUES(?,?,?,?,?)`, audit.RecordID, audit.Action, audit.Operator, audit.OccurredAt, audit.Detail)
	return err
}

func (s *Store) ListAudits(ctx context.Context, recordID int64) ([]model.AdminAudit, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,record_id,action,operator,occurred_at,detail FROM admin_audits WHERE record_id=? ORDER BY id`, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.AdminAudit{}
	for rows.Next() {
		var a model.AdminAudit
		if err := rows.Scan(&a.ID, &a.RecordID, &a.Action, &a.Operator, &a.OccurredAt, &a.Detail); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
