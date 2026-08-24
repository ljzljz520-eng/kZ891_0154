package model

import (
	"fmt"
	"sort"
	"strings"
)

type Change struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func DiffWarranty(before, after WarrantyRecord) []Change {
	changes := []Change{}
	if before.Phone != after.Phone {
		changes = append(changes, Change{"phone", before.Phone, after.Phone})
	}
	if before.SerialNumber != after.SerialNumber {
		changes = append(changes, Change{"serial_number", before.SerialNumber, after.SerialNumber})
	}
	if before.ExpiryDate != after.ExpiryDate {
		changes = append(changes, Change{"expiry_date", before.ExpiryDate, after.ExpiryDate})
	}
	if before.PurchaseChannel != after.PurchaseChannel {
		changes = append(changes, Change{"purchase_channel", before.PurchaseChannel, after.PurchaseChannel})
	}
	if before.Category != after.Category {
		changes = append(changes, Change{"category", before.Category, after.Category})
	}
	if before.UpdatedAt != after.UpdatedAt {
		changes = append(changes, Change{"updated_at", before.UpdatedAt, after.UpdatedAt})
	}
	if before.Active != after.Active {
		changes = append(changes, Change{"active", fmt.Sprint(before.Active), fmt.Sprint(after.Active)})
	}
	if strings.Join(before.ServicePoints, "|") != strings.Join(after.ServicePoints, "|") {
		changes = append(changes, Change{"service_points", strings.Join(before.ServicePoints, ","), strings.Join(after.ServicePoints, ",")})
	}
	return changes
}

func ChangeFields(changes []Change) []string {
	result := make([]string, 0, len(changes))
	for _, change := range changes {
		result = append(result, change.Field)
	}
	sort.Strings(result)
	return result
}
func ChangeText(changes []Change) string {
	parts := make([]string, 0, len(changes))
	for _, change := range changes {
		parts = append(parts, change.Field+":"+change.Before+"->"+change.After)
	}
	return strings.Join(parts, "; ")
}
func SameIdentity(a, b WarrantyRecord) bool {
	return NormalizePhone(a.Phone) == NormalizePhone(b.Phone) && NormalizeSerial(a.SerialNumber) == NormalizeSerial(b.SerialNumber)
}
func SamePoints(a, b []string) bool {
	left, right := append([]string(nil), a...), append([]string(nil), b...)
	sort.Strings(left)
	sort.Strings(right)
	return strings.Join(left, "|") == strings.Join(right, "|")
}
func MergeLabels(values ...[]string) []string {
	set := map[string]bool{}
	for _, group := range values {
		for _, value := range group {
			if strings.TrimSpace(value) != "" {
				set[value] = true
			}
		}
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
