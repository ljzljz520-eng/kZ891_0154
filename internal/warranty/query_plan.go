package warranty

import (
	"fmt"
	"strings"

	"warrantyservice/internal/model"
)

type QueryPlan struct {
	NormalizedPhone  string   `json:"normalized_phone"`
	NormalizedSerial string   `json:"normalized_serial"`
	Checks           []string `json:"checks"`
	Warnings         []string `json:"warnings"`
}

func PlanQuery(phone, serial string) QueryPlan {
	plan := QueryPlan{NormalizedPhone: model.NormalizePhone(phone), NormalizedSerial: model.NormalizeSerial(serial), Checks: []string{}, Warnings: []string{}}
	if plan.NormalizedPhone == "" {
		plan.Warnings = append(plan.Warnings, "手机号为空")
	} else {
		plan.Checks = append(plan.Checks, "手机号格式")
	}
	if plan.NormalizedSerial == "" {
		plan.Warnings = append(plan.Warnings, "序列号为空")
	} else {
		plan.Checks = append(plan.Checks, "序列号长度")
	}
	if strings.Contains(plan.NormalizedSerial, "-") {
		plan.Checks = append(plan.Checks, "序列号前缀")
	}
	return plan
}

func (p QueryPlan) Ready() bool { return len(p.Warnings) == 0 && len(p.Checks) >= 2 }
func (p QueryPlan) Description() string {
	return fmt.Sprintf("手机号=%s，序列号=%s，检查=%d，警告=%d", p.NormalizedPhone, p.NormalizedSerial, len(p.Checks), len(p.Warnings))
}
func (p QueryPlan) HasCheck(value string) bool {
	for _, check := range p.Checks {
		if check == value {
			return true
		}
	}
	return false
}
func (p QueryPlan) HasWarning(value string) bool {
	for _, warning := range p.Warnings {
		if warning == value {
			return true
		}
	}
	return false
}

func QueryResultLabel(result model.QueryResult) string {
	switch result.Status {
	case model.StatusValid:
		return "在保"
	case model.StatusExpired:
		return "已过期"
	case model.StatusMissing:
		return "未找到"
	case model.StatusInvalid:
		return "参数错误"
	default:
		return "未知"
	}
}

func ExplainResult(result model.QueryResult) []string {
	lines := []string{QueryResultLabel(result)}
	if result.Record != nil {
		lines = append(lines, "设备序列号："+result.Record.SerialNumber, "适用品类："+result.Record.DisplayCategory(), "购买渠道："+result.Record.PurchaseChannel, "到期日："+result.Record.ExpiryDate)
	}
	if result.Advice != nil {
		lines = append(lines, result.Advice.Message, "客服热线："+result.Advice.Hotline)
	}
	if result.Message != "" {
		lines = append(lines, result.Message)
	}
	return lines
}

func RequireRecord(result model.QueryResult) (*model.WarrantyRecord, error) {
	if result.Record == nil {
		return nil, fmt.Errorf("warranty record is missing")
	}
	return result.Record, nil
}
