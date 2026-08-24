package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"warrantyservice/internal/config"
	"warrantyservice/internal/warranty"
)

func ValidateAdminRequest(r *http.Request, cfg config.Config) error {
	if r == nil {
		return errors.New("request is nil")
	}
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return errors.New("method is not allowed")
	}
	if !AdminTokenValid(r, cfg.AdminToken) {
		return errors.New("管理员令牌无效")
	}
	if strings.TrimSpace(r.Header.Get("X-Operator")) == "" {
		return errors.New("管理员身份不能为空")
	}
	return nil
}

func ValidateAdminInput(input warranty.UpsertInput, policy config.BusinessPolicy) error {
	input = warranty.NormalizeInput(input)
	if err := warranty.ValidateUpsertInput(input); err != nil {
		return err
	}
	if err := policy.ValidateCategory(input.Category); err != nil {
		return err
	}
	if err := policy.ValidateChannel(input.PurchaseChannel); err != nil {
		return err
	}
	return policy.ValidatePointCount(len(input.ServicePoints))
}

func AdminActionLabel(action string) string {
	switch action {
	case "create":
		return "新增"
	case "update":
		return "修改"
	case "deactivate":
		return "停用"
	default:
		return "其他"
	}
}
