package config

import (
	"errors"
	"sort"
	"strings"
)

type BusinessPolicy struct {
	AllowedCategories []string
	AllowedChannels   []string
	MaximumPoints     int
	RequireOperator   bool
}

func DefaultPolicy() BusinessPolicy {
	return BusinessPolicy{AllowedCategories: []string{"电视", "冰箱", "洗衣机", "空调", "热水器", "厨房电器"}, AllowedChannels: []string{"官方商城", "授权门店", "电商旗舰店", "企业团购"}, MaximumPoints: 8, RequireOperator: true}
}

func (p BusinessPolicy) ValidateCategory(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("category is required")
	}
	for _, candidate := range p.AllowedCategories {
		if candidate == value {
			return nil
		}
	}
	return errors.New("category is not supported")
}

func (p BusinessPolicy) ValidateChannel(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("purchase channel is required")
	}
	for _, candidate := range p.AllowedChannels {
		if candidate == value {
			return nil
		}
	}
	return errors.New("purchase channel is not supported")
}

func (p BusinessPolicy) ValidatePointCount(count int) error {
	if count < 0 {
		return errors.New("point count cannot be negative")
	}
	if p.MaximumPoints > 0 && count > p.MaximumPoints {
		return errors.New("too many service points")
	}
	return nil
}

func (p BusinessPolicy) Categories() []string {
	result := append([]string(nil), p.AllowedCategories...)
	sort.Strings(result)
	return result
}
func (p BusinessPolicy) Channels() []string {
	result := append([]string(nil), p.AllowedChannels...)
	sort.Strings(result)
	return result
}
func (p BusinessPolicy) SupportsCategory(value string) bool { return p.ValidateCategory(value) == nil }
func (p BusinessPolicy) SupportsChannel(value string) bool  { return p.ValidateChannel(value) == nil }
