package model

import (
	"fmt"
	"sort"
	"strings"
)

type WarrantySummary struct {
	Key             string `json:"key"`
	SerialNumber    string `json:"serial_number"`
	Category        string `json:"category"`
	PurchaseChannel string `json:"purchase_channel"`
	ExpiryDate      string `json:"expiry_date"`
	PointCount      int    `json:"point_count"`
	Active          bool   `json:"active"`
}

func (r WarrantyRecord) Summary() WarrantySummary {
	return WarrantySummary{Key: r.Key(), SerialNumber: r.SerialNumber, Category: r.DisplayCategory(), PurchaseChannel: r.PurchaseChannel, ExpiryDate: r.ExpiryDate, PointCount: len(r.ServicePoints), Active: r.Active}
}

func (r WarrantyRecord) Clone() WarrantyRecord {
	copy := r
	copy.ServicePoints = append([]string(nil), r.ServicePoints...)
	return copy
}

func (r WarrantyRecord) Validate() error {
	if NormalizePhone(r.Phone) != r.Phone {
		return fmt.Errorf("phone is not normalized")
	}
	if NormalizeSerial(r.SerialNumber) != r.SerialNumber {
		return fmt.Errorf("serial is not normalized")
	}
	if r.ID <= 0 {
		return fmt.Errorf("id must be positive")
	}
	if len(r.ExpiryDate) != 10 {
		return fmt.Errorf("expiry date must be YYYY-MM-DD")
	}
	if r.Category == "" {
		return fmt.Errorf("category is required")
	}
	if r.PurchaseChannel == "" {
		return fmt.Errorf("purchase channel is required")
	}
	return nil
}

func (r WarrantyRecord) Labels() []string {
	labels := []string{r.DisplayCategory(), r.PurchaseChannel}
	if r.Active {
		labels = append(labels, "在保")
	} else {
		labels = append(labels, "停用")
	}
	if len(r.ServicePoints) > 0 {
		labels = append(labels, fmt.Sprintf("%d个服务网点", len(r.ServicePoints)))
	}
	return labels
}

func SortRecords(records []WarrantyRecord) []WarrantyRecord {
	result := append([]WarrantyRecord(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Category == result[j].Category {
			return result[i].SerialNumber < result[j].SerialNumber
		}
		return result[i].Category < result[j].Category
	})
	return result
}

func SearchRecords(records []WarrantyRecord, query string) []WarrantyRecord {
	query = strings.ToUpper(strings.TrimSpace(query))
	if query == "" {
		return append([]WarrantyRecord(nil), records...)
	}
	result := make([]WarrantyRecord, 0)
	for _, record := range records {
		if strings.Contains(strings.ToUpper(record.SerialNumber), query) || strings.Contains(record.Category, query) || strings.Contains(record.Phone, query) {
			result = append(result, record)
		}
	}
	return result
}

func PointNames(points []ServicePoint) []string {
	result := make([]string, 0, len(points))
	for _, point := range points {
		if point.Name != "" {
			result = append(result, point.Name)
		}
	}
	sort.Strings(result)
	return result
}

func (p ServicePoint) Validate() error {
	if p.ID <= 0 {
		return fmt.Errorf("point id must be positive")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("point name is required")
	}
	if strings.TrimSpace(p.Address) == "" {
		return fmt.Errorf("point address is required")
	}
	if strings.TrimSpace(p.Phone) == "" {
		return fmt.Errorf("point phone is required")
	}
	return nil
}

func (a AdminAudit) Summary() string {
	return fmt.Sprintf("%s by %s on %s", a.Action, a.Operator, a.OccurredAt)
}
func (a AdminAudit) IsMutation() bool {
	return a.Action == "create" || a.Action == "update" || a.Action == "deactivate"
}
func (a SupportAdvice) ContactLine() string { return fmt.Sprintf("%s（%s）", a.Hotline, a.Channels) }
func (a SupportAdvice) Words() []string {
	return strings.Fields(strings.ReplaceAll(a.Message, "，", " "))
}
