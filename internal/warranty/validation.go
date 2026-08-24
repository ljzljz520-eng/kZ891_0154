package warranty

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"warrantyservice/internal/model"
)

var datePattern = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)

func ValidateUpsertInput(input UpsertInput) error {
	if input.ID <= 0 {
		return errors.New("记录ID必须为正数")
	}
	phone, serial := model.NormalizePhone(input.Phone), model.NormalizeSerial(input.SerialNumber)
	if err := ValidateIdentity(phone, serial); err != nil {
		return err
	}
	if !datePattern.MatchString(input.ExpiryDate) {
		return errors.New("到期日必须使用YYYY-MM-DD")
	}
	if input.UpdatedAt != "" && !datePattern.MatchString(input.UpdatedAt) {
		return errors.New("更新时间必须使用YYYY-MM-DD")
	}
	if strings.TrimSpace(input.Category) == "" {
		return errors.New("适用品类不能为空")
	}
	if strings.TrimSpace(input.PurchaseChannel) == "" {
		return errors.New("购买渠道不能为空")
	}
	seen := map[int64]bool{}
	for _, point := range input.ServicePoints {
		if seen[point.ID] {
			return fmt.Errorf("服务网点重复：%d", point.ID)
		}
		seen[point.ID] = true
		if err := point.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ValidateExpiry(today, expiry string) error {
	if !datePattern.MatchString(today) || !datePattern.MatchString(expiry) {
		return errors.New("日期格式无效")
	}
	if expiry < today {
		return errors.New("到期日不能早于当前业务日期")
	}
	return nil
}

func NormalizeInput(input UpsertInput) UpsertInput {
	input.Phone = model.NormalizePhone(input.Phone)
	input.SerialNumber = model.NormalizeSerial(input.SerialNumber)
	input.Category = strings.TrimSpace(input.Category)
	input.PurchaseChannel = strings.TrimSpace(input.PurchaseChannel)
	for i := range input.ServicePoints {
		input.ServicePoints[i].Name = strings.TrimSpace(input.ServicePoints[i].Name)
		input.ServicePoints[i].Address = strings.TrimSpace(input.ServicePoints[i].Address)
	}
	return input
}

func ActiveCount(records []model.WarrantyRecord, today string) int {
	count := 0
	for _, record := range records {
		if record.Active && record.ExpiryDate >= today {
			count++
		}
	}
	return count
}
func ExpiredCount(records []model.WarrantyRecord, today string) int {
	count := 0
	for _, record := range records {
		if !record.Active || record.ExpiryDate < today {
			count++
		}
	}
	return count
}

func Categories(records []model.WarrantyRecord) []string {
	set := map[string]bool{}
	for _, record := range records {
		if record.Category != "" {
			set[record.Category] = true
		}
	}
	result := make([]string, 0, len(set))
	for category := range set {
		result = append(result, category)
	}
	sort.Strings(result)
	return result
}

func Channels(records []model.WarrantyRecord) []string {
	set := map[string]bool{}
	for _, record := range records {
		if record.PurchaseChannel != "" {
			set[record.PurchaseChannel] = true
		}
	}
	result := make([]string, 0, len(set))
	for channel := range set {
		result = append(result, channel)
	}
	sort.Strings(result)
	return result
}

func GroupByCategory(records []model.WarrantyRecord) map[string][]model.WarrantyRecord {
	groups := map[string][]model.WarrantyRecord{}
	for _, record := range records {
		groups[record.Category] = append(groups[record.Category], record)
	}
	return groups
}

func Summaries(records []model.WarrantyRecord) []model.WarrantySummary {
	result := make([]model.WarrantySummary, 0, len(records))
	for _, record := range records {
		result = append(result, record.Summary())
	}
	return result
}

func (s *Service) ValidateForQuery(phone, serial string) (string, string, error) {
	phone = model.NormalizePhone(phone)
	serial = model.NormalizeSerial(serial)
	if err := ValidateIdentity(phone, serial); err != nil {
		return "", "", err
	}
	return phone, serial, nil
}
