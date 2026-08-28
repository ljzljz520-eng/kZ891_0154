package model

import "strings"

type WarrantyStatus string

const (
	StatusValid   WarrantyStatus = "valid"
	StatusExpired WarrantyStatus = "expired"
	StatusMissing WarrantyStatus = "missing"
	StatusInvalid WarrantyStatus = "invalid"
)

type QueryResult struct {
	Status        WarrantyStatus  `json:"status"`
	Record        *WarrantyRecord `json:"record,omitempty"`
	Advice        *SupportAdvice  `json:"advice,omitempty"`
	ServicePoints []ServicePoint  `json:"service_points,omitempty"`
	Message       string          `json:"message,omitempty"`
}

func (s WarrantyStatus) IsSuccess() bool { return s == StatusValid || s == StatusExpired }
func (s WarrantyStatus) HTTPCode() int {
	if s == StatusMissing {
		return 404
	}
	if s == StatusInvalid {
		return 400
	}
	return 200
}
func (s WarrantyStatus) String() string { return strings.ToLower(string(s)) }
func NewInvalidResult(message string) QueryResult {
	return QueryResult{Status: StatusInvalid, Message: message}
}
func NewMissingResult(advice *SupportAdvice) QueryResult {
	return QueryResult{Status: StatusMissing, Advice: advice, Message: "未找到延保记录"}
}
