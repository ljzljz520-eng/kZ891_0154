package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
	"warrantyservice/internal/model"
)

type Store struct {
	db   *sql.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if path != ":memory:" && path != "file::memory:?cache=shared" {
		if dir := directory(path); dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, err
			}
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db, path: path}
	if err := s.initialize(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func directory(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i == 0 {
				return "/"
			}
			return path[:i]
		}
	}
	return ""
}

func (s *Store) initialize(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS warranties (id INTEGER PRIMARY KEY, phone TEXT NOT NULL, serial_number TEXT NOT NULL, expiry_date TEXT NOT NULL, purchase_channel TEXT NOT NULL, category TEXT NOT NULL, updated_at TEXT NOT NULL, active INTEGER NOT NULL DEFAULT 1, UNIQUE(phone, serial_number))`,
		`CREATE TABLE IF NOT EXISTS service_points (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE, address TEXT NOT NULL, phone TEXT NOT NULL, categories TEXT NOT NULL, open_hours TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS warranty_points (warranty_id INTEGER NOT NULL, point_id INTEGER NOT NULL, PRIMARY KEY(warranty_id, point_id))`,
		`CREATE TABLE IF NOT EXISTS admin_audits (id INTEGER PRIMARY KEY AUTOINCREMENT, record_id INTEGER NOT NULL, action TEXT NOT NULL, operator TEXT NOT NULL, occurred_at TEXT NOT NULL, detail TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS support_advices (id INTEGER PRIMARY KEY AUTOINCREMENT, phone TEXT NOT NULL, serial_number TEXT NOT NULL, message TEXT NOT NULL, created_at TEXT NOT NULL, hotline TEXT NOT NULL, channels TEXT NOT NULL)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("initialize schema: %w", err)
		}
	}
	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func (s *Store) Reopen() error {
	if s == nil {
		return errors.New("nil store")
	}
	if err := s.Close(); err != nil {
		return err
	}
	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	s.db = db
	return s.initialize(context.Background())
}

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) SaveWarranty(ctx context.Context, record model.WarrantyRecord, points []model.ServicePoint) error {
	if record.ID == 0 {
		return errors.New("warranty id is required")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	rollback := func(e error) error { _ = tx.Rollback(); return e }
	_, err = tx.ExecContext(ctx, `INSERT INTO warranties(id,phone,serial_number,expiry_date,purchase_channel,category,updated_at,active) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(phone,serial_number) DO UPDATE SET id=excluded.id, expiry_date=excluded.expiry_date, purchase_channel=excluded.purchase_channel, category=excluded.category, updated_at=excluded.updated_at, active=excluded.active`, record.ID, record.Phone, record.SerialNumber, record.ExpiryDate, record.PurchaseChannel, record.Category, record.UpdatedAt, boolInt(record.Active))
	if err != nil {
		return rollback(fmt.Errorf("save warranty: %w", err))
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM warranty_points WHERE warranty_id=?`, record.ID); err != nil {
		return rollback(err)
	}
	for _, point := range points {
		if err := savePointTx(ctx, tx, point); err != nil {
			return rollback(err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO warranty_points(warranty_id,point_id) VALUES(?,?)`, record.ID, point.ID); err != nil {
			return rollback(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func savePointTx(ctx context.Context, tx *sql.Tx, point model.ServicePoint) error {
	if point.ID == 0 || point.Name == "" {
		return errors.New("service point identity is required")
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO service_points(id,name,address,phone,categories,open_hours) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,address=excluded.address,phone=excluded.phone,categories=excluded.categories,open_hours=excluded.open_hours`, point.ID, point.Name, point.Address, point.Phone, join(point.Categories), point.OpenHours)
	return err
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
func intBool(v int) bool { return v != 0 }
func join(values []string) string {
	out := ""
	for i, v := range values {
		if i > 0 {
			out += "|"
		}
		out += v
	}
	return out
}
func split(value string) []string {
	if value == "" {
		return nil
	}
	out := []string{}
	start := 0
	for i := 0; i <= len(value); i++ {
		if i == len(value) || value[i] == '|' {
			out = append(out, value[start:i])
			start = i + 1
		}
	}
	return out
}
