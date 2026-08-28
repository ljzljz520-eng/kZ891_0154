package support

import (
	"html/template"
	"strings"

	"warrantyservice/internal/model"
)

func RenderAdviceText(advice *model.SupportAdvice) string {
	if advice == nil {
		return "请联系客服核对延保信息。"
	}
	parts := []string{advice.Message}
	if advice.Hotline != "" {
		parts = append(parts, "热线："+advice.Hotline)
	}
	if advice.Channels != "" {
		parts = append(parts, "服务渠道："+advice.Channels)
	}
	return strings.Join(parts, " ")
}

func SafeAdviceHTML(advice *model.SupportAdvice) template.HTML {
	return template.HTML(template.HTMLEscapeString(RenderAdviceText(advice)))
}

func NormalizeChannels(v string) string {
	if strings.TrimSpace(v) == "" {
		return "电话、官方小程序、授权服务网点"
	}
	return strings.TrimSpace(v)
}
