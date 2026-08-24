package support

import (
	"errors"
	"fmt"
	"strings"

	"warrantyservice/internal/model"
)

type Escalation struct {
	Level      string `json:"level"`
	Reason     string `json:"reason"`
	Owner      string `json:"owner"`
	NextAction string `json:"next_action"`
}

type Case struct {
	Phone      string              `json:"phone"`
	Serial     string              `json:"serial"`
	Category   string              `json:"category"`
	Advice     model.SupportAdvice `json:"advice"`
	Escalation Escalation          `json:"escalation"`
}

func ClassifyMissing(phone, serial string) Escalation {
	if phone == "" || serial == "" {
		return Escalation{Level: "L1", Reason: "查询信息不完整", Owner: "在线客服", NextAction: "补充手机号和设备序列号"}
	}
	if strings.HasPrefix(serial, "TV-") || strings.HasPrefix(serial, "AC-") {
		return Escalation{Level: "L2", Reason: "高价值家电需人工核对", Owner: "区域服务中心", NextAction: "转接区域服务中心"}
	}
	return Escalation{Level: "L1", Reason: "系统中未找到记录", Owner: "在线客服", NextAction: "核对购买凭证"}
}

func NewCase(advice model.SupportAdvice) Case {
	return Case{Phone: advice.Phone, Serial: advice.SerialNumber, Advice: advice, Escalation: ClassifyMissing(advice.Phone, advice.SerialNumber)}
}

func ValidateCase(c Case) error {
	if c.Phone == "" {
		return errors.New("case phone is required")
	}
	if c.Serial == "" {
		return errors.New("case serial is required")
	}
	if err := ValidateAdviceFields(c.Advice); err != nil {
		return err
	}
	if c.Escalation.Level == "" {
		return errors.New("case escalation is missing")
	}
	return nil
}

func ValidateAdviceFields(a model.SupportAdvice) error {
	if a.Phone == "" || a.SerialNumber == "" {
		return errors.New("advice identity is required")
	}
	if a.Message == "" || a.Hotline == "" {
		return errors.New("advice contact is required")
	}
	return nil
}

func EscalationText(e Escalation) string {
	return fmt.Sprintf("%s：%s；负责人：%s；下一步：%s", e.Level, e.Reason, e.Owner, e.NextAction)
}

func BuildCaseSummary(c Case) string {
	if c.Escalation.Level == "" {
		return "待补充客服工单信息"
	}
	return fmt.Sprintf("设备%s未找到。%s", c.Serial, EscalationText(c.Escalation))
}

func ContactOptions(hotline string) []string {
	if strings.TrimSpace(hotline) == "" {
		hotline = "400-800-1234"
	}
	return []string{"电话 " + hotline, "官方小程序", "授权服务网点"}
}
