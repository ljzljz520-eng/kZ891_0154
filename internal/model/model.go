package model

import "strings"

type WarrantyRecord struct {
	ID              int64    `json:"id"`
	Phone           string   `json:"phone"`
	SerialNumber    string   `json:"serial_number"`
	ExpiryDate      string   `json:"expiry_date"`
	ServicePoints   []string `json:"service_points"`
	PurchaseChannel string   `json:"purchase_channel"`
	Category        string   `json:"category"`
	UpdatedAt       string   `json:"updated_at"`
	Active          bool     `json:"active"`
}

type ServicePoint struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	Phone      string   `json:"phone"`
	Categories []string `json:"categories"`
	OpenHours  string   `json:"open_hours"`
}

type AdminAudit struct {
	ID         int64  `json:"id"`
	RecordID   int64  `json:"record_id"`
	Action     string `json:"action"`
	Operator   string `json:"operator"`
	OccurredAt string `json:"occurred_at"`
	Detail     string `json:"detail"`
}

type SupportAdvice struct {
	ID           int64  `json:"id"`
	Phone        string `json:"phone"`
	SerialNumber string `json:"serial_number"`
	Message      string `json:"message"`
	CreatedAt    string `json:"created_at"`
	Hotline      string `json:"hotline"`
	Channels     string `json:"channels"`
}

func NormalizePhone(v string) string  { return strings.TrimSpace(v) }
func NormalizeSerial(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
func (r WarrantyRecord) Key() string {
	return NormalizePhone(r.Phone) + ":" + NormalizeSerial(r.SerialNumber)
}
func (r WarrantyRecord) IsActive() bool { return r.Active }
func (r WarrantyRecord) DisplayCategory() string {
	if r.Category == "" {
		return "未分类"
	}
	return r.Category
}
func (p ServicePoint) DisplayName() string {
	if p.Name == "" {
		return "未命名网点"
	}
	return p.Name
}
func (a SupportAdvice) IsActionable() bool { return a.Message != "" && a.Hotline != "" }
