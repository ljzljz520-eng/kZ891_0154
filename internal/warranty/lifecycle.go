package warranty

import (
	"fmt"
	"strings"

	"warrantyservice/internal/model"
)

type LifecycleEvent struct {
	Name    string
	Detail  string
	Allowed bool
}

func EventsFor(record model.WarrantyRecord, today string) []LifecycleEvent {
	events := []LifecycleEvent{{Name: "created", Detail: record.UpdatedAt, Allowed: record.ID > 0}}
	if record.Active {
		events = append(events, LifecycleEvent{Name: "active", Detail: "可查询", Allowed: record.ExpiryDate >= today})
	} else {
		events = append(events, LifecycleEvent{Name: "inactive", Detail: "已停用", Allowed: false})
	}
	if record.ExpiryDate < today {
		events = append(events, LifecycleEvent{Name: "expired", Detail: record.ExpiryDate, Allowed: false})
	}
	return events
}

func CanUpdate(record model.WarrantyRecord, operator string) error {
	if strings.TrimSpace(operator) == "" {
		return fmt.Errorf("管理员不能为空")
	}
	if record.ID <= 0 {
		return fmt.Errorf("记录ID无效")
	}
	return nil
}

func ActionFor(existing *model.WarrantyRecord, active bool) string {
	if existing == nil {
		return "create"
	}
	if existing.Active && !active {
		return "deactivate"
	}
	return "update"
}

func MergeRecord(old model.WarrantyRecord, input UpsertInput) model.WarrantyRecord {
	merged := old.Clone()
	if input.Phone != "" {
		merged.Phone = model.NormalizePhone(input.Phone)
	}
	if input.SerialNumber != "" {
		merged.SerialNumber = model.NormalizeSerial(input.SerialNumber)
	}
	if input.ExpiryDate != "" {
		merged.ExpiryDate = input.ExpiryDate
	}
	if input.Category != "" {
		merged.Category = input.Category
	}
	if input.PurchaseChannel != "" {
		merged.PurchaseChannel = input.PurchaseChannel
	}
	if input.UpdatedAt != "" {
		merged.UpdatedAt = input.UpdatedAt
	}
	merged.Active = input.Active
	if input.ServicePoints != nil {
		merged.ServicePoints = nil
		for _, point := range input.ServicePoints {
			merged.ServicePoints = append(merged.ServicePoints, point.Name)
		}
	}
	return merged
}
