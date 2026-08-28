package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"warrantyservice/internal/model"
)

type Snapshot struct {
	Warranties []model.WarrantyRecord `json:"warranties"`
	Points     []model.ServicePoint   `json:"service_points"`
	Audits     []model.AdminAudit     `json:"audits"`
	Advices    []model.SupportAdvice  `json:"advices"`
}

func (s *Store) Export(ctx context.Context) (Snapshot, error) {
	warranties, err := s.ListWarranties(ctx, false)
	if err != nil {
		return Snapshot{}, err
	}
	points, err := s.allPoints(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{Warranties: warranties, Points: points}
	for _, record := range warranties {
		audits, auditErr := s.ListAudits(ctx, record.ID)
		if auditErr != nil {
			return Snapshot{}, auditErr
		}
		snapshot.Audits = append(snapshot.Audits, audits...)
	}
	return snapshot, nil
}

func (s *Store) allPoints(ctx context.Context) ([]model.ServicePoint, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,name,address,phone,categories,open_hours FROM service_points ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.ServicePoint{}
	for rows.Next() {
		var p model.ServicePoint
		var categories string
		if err := rows.Scan(&p.ID, &p.Name, &p.Address, &p.Phone, &categories, &p.OpenHours); err != nil {
			return nil, err
		}
		p.Categories = split(categories)
		result = append(result, p)
	}
	return result, rows.Err()
}

func EncodeSnapshot(snapshot Snapshot) ([]byte, error) { return json.MarshalIndent(snapshot, "", "  ") }
func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(strings.TrimSpace(string(data))) == 0 {
		return Snapshot{}, errors.New("snapshot is empty")
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("decode snapshot: %w", err)
	}
	return snapshot, nil
}

func (s *Store) Import(ctx context.Context, snapshot Snapshot) error {
	if len(snapshot.Warranties) == 0 {
		return errors.New("snapshot has no warranties")
	}
	for _, point := range snapshot.Points {
		if err := s.UpsertServicePoint(ctx, point); err != nil {
			return err
		}
	}
	for _, record := range snapshot.Warranties {
		points, err := s.FindServicePoints(ctx, record.ServicePoints)
		if err != nil {
			return err
		}
		if err := s.SaveWarranty(ctx, record, points); err != nil {
			return err
		}
	}
	return nil
}

func SnapshotKeys(snapshot Snapshot) []string {
	keys := make([]string, 0, len(snapshot.Warranties))
	for _, record := range snapshot.Warranties {
		keys = append(keys, record.Key())
	}
	return keys
}
