package httpapi

import (
	"net/http"
	"sort"
	"strings"

	"warrantyservice/internal/model"
)

type QueryFilter struct {
	Phone      string
	Serial     string
	Category   string
	Channel    string
	ActiveOnly bool
	Page       int
	PageSize   int
}

func NewQueryFilter(r *http.Request) QueryFilter {
	filter := QueryFilter{Page: 1, PageSize: 20}
	if r == nil {
		return filter
	}
	query := r.URL.Query()
	filter.Phone = strings.TrimSpace(query.Get("phone"))
	filter.Serial = strings.TrimSpace(query.Get("serial"))
	filter.Category = strings.TrimSpace(query.Get("category"))
	filter.Channel = strings.TrimSpace(query.Get("channel"))
	filter.ActiveOnly = ParseBool(query.Get("active"), false)
	return filter
}

func FilterRecords(records []model.WarrantyRecord, filter QueryFilter, today string) []model.WarrantyRecord {
	result := make([]model.WarrantyRecord, 0, len(records))
	for _, record := range records {
		if filter.Phone != "" && record.Phone != filter.Phone {
			continue
		}
		if filter.Serial != "" && !strings.Contains(strings.ToUpper(record.SerialNumber), strings.ToUpper(filter.Serial)) {
			continue
		}
		if filter.Category != "" && record.Category != filter.Category {
			continue
		}
		if filter.Channel != "" && record.PurchaseChannel != filter.Channel {
			continue
		}
		if filter.ActiveOnly && (!record.Active || record.ExpiryDate < today) {
			continue
		}
		result = append(result, record)
	}
	return result
}

func PaginateRecords(records []model.WarrantyRecord, page, size int) []model.WarrantyRecord {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= len(records) {
		return []model.WarrantyRecord{}
	}
	end := start + size
	if end > len(records) {
		end = len(records)
	}
	return append([]model.WarrantyRecord(nil), records[start:end]...)
}

func SortForDisplay(records []model.WarrantyRecord, by string) []model.WarrantyRecord {
	result := append([]model.WarrantyRecord(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		switch by {
		case "expiry":
			return result[i].ExpiryDate < result[j].ExpiryDate
		case "category":
			return result[i].Category < result[j].Category
		default:
			return result[i].ID < result[j].ID
		}
	})
	return result
}

func FilterSummary(records []model.WarrantyRecord) map[string]int {
	result := map[string]int{"total": len(records), "active": 0, "inactive": 0}
	for _, record := range records {
		if record.Active {
			result["active"]++
		} else {
			result["inactive"]++
		}
	}
	return result
}
