package warranty

import (
	"context"
	"fmt"
	"sort"

	"warrantyservice/internal/model"
)

type Report struct {
	Total      int            `json:"total"`
	Active     int            `json:"active"`
	Expired    int            `json:"expired"`
	ByCategory map[string]int `json:"by_category"`
	ByChannel  map[string]int `json:"by_channel"`
}

type WarrantyLister interface {
	ListWarranties(context.Context, bool) ([]model.WarrantyRecord, error)
}

func BuildReport(ctx context.Context, lister WarrantyLister, today string) (Report, error) {
	if lister == nil {
		return Report{}, fmt.Errorf("warranty lister is unavailable")
	}
	records, err := lister.ListWarranties(ctx, false)
	if err != nil {
		return Report{}, err
	}
	report := Report{ByCategory: map[string]int{}, ByChannel: map[string]int{}}
	for _, record := range records {
		report.Total++
		if record.Active && record.ExpiryDate >= today {
			report.Active++
		} else {
			report.Expired++
		}
		report.ByCategory[record.Category]++
		report.ByChannel[record.PurchaseChannel]++
	}
	return report, nil
}

func SortedCategories(report Report) []string {
	result := make([]string, 0, len(report.ByCategory))
	for key := range report.ByCategory {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
func SortedChannels(report Report) []string {
	result := make([]string, 0, len(report.ByChannel))
	for key := range report.ByChannel {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
