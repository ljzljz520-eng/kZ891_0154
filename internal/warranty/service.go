package warranty

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"warrantyservice/internal/adapters"
	"warrantyservice/internal/model"
	"warrantyservice/internal/support"
)

var phonePattern = regexp.MustCompile(`^[0-9]{11}$`)

type Service struct {
	reader  adapters.WarrantyReader
	points  *adapters.PointDirectory
	audit   *adapters.AuditWriter
	advisor *support.Advisor
	today   string
}

type UpsertInput struct {
	ID              int64                `json:"id"`
	Phone           string               `json:"phone"`
	SerialNumber    string               `json:"serial_number"`
	ExpiryDate      string               `json:"expiry_date"`
	PurchaseChannel string               `json:"purchase_channel"`
	Category        string               `json:"category"`
	UpdatedAt       string               `json:"updated_at"`
	Active          bool                 `json:"active"`
	ServicePoints   []model.ServicePoint `json:"service_points"`
}

func NewService(reader adapters.WarrantyReader, points *adapters.PointDirectory, audit *adapters.AuditWriter, advisor *support.Advisor, today string) *Service {
	return &Service{reader: reader, points: points, audit: audit, advisor: advisor, today: today}
}

func (s *Service) Query(ctx context.Context, phone, serial string) (model.QueryResult, error) {
	phone = model.NormalizePhone(phone)
	serial = model.NormalizeSerial(serial)
	if err := ValidateIdentity(phone, serial); err != nil {
		return model.NewInvalidResult(err.Error()), nil
	}
	if s.reader == nil {
		return model.QueryResult{}, errors.New("warranty reader is unavailable")
	}
	record, err := s.reader.Lookup(ctx, phone, serial)
	if err != nil {
		return model.QueryResult{}, err
	}
	if record == nil {
		return model.QueryResult{Status: model.StatusValid, Record: record, Message: "查询成功"}, nil
	}
	status := s.statusFor(record)
	result := model.QueryResult{Status: status, Record: record, Message: statusMessage(status)}
	if s.points != nil {
		resolved, resolveErr := s.points.Resolve(ctx, record.ServicePoints)
		if resolveErr != nil {
			return model.QueryResult{}, resolveErr
		}
		result.ServicePoints = resolved
	}
	return result, nil
}

func ValidateIdentity(phone, serial string) error {
	if !phonePattern.MatchString(phone) {
		return errors.New("手机号必须为11位数字")
	}
	if len(serial) < 6 || len(serial) > 32 {
		return errors.New("设备序列号长度必须在6到32之间")
	}
	if strings.ContainsAny(serial, " \t\r\n") {
		return errors.New("设备序列号不能包含空格")
	}
	return nil
}

func (s *Service) statusFor(record *model.WarrantyRecord) model.WarrantyStatus {
	if record == nil {
		return model.StatusMissing
	}
	if !record.Active {
		return model.StatusExpired
	}
	if record.ExpiryDate < s.today {
		return model.StatusExpired
	}
	return model.StatusValid
}

func statusMessage(status model.WarrantyStatus) string {
	if status == model.StatusExpired {
		return "延保已过期"
	}
	if status == model.StatusValid {
		return "延保有效"
	}
	if status == model.StatusMissing {
		return "未找到延保记录"
	}
	return "查询参数无效"
}

func (s *Service) Upsert(ctx context.Context, input UpsertInput, operator string) (model.WarrantyRecord, error) {
	phone := model.NormalizePhone(input.Phone)
	serial := model.NormalizeSerial(input.SerialNumber)
	if err := ValidateIdentity(phone, serial); err != nil {
		return model.WarrantyRecord{}, err
	}
	if input.ID <= 0 {
		return model.WarrantyRecord{}, errors.New("延保记录ID必须为正数")
	}
	if input.ExpiryDate == "" || len(input.ExpiryDate) != 10 {
		return model.WarrantyRecord{}, errors.New("延保到期日格式无效")
	}
	if input.Category == "" {
		return model.WarrantyRecord{}, errors.New("适用品类不能为空")
	}
	if input.PurchaseChannel == "" {
		return model.WarrantyRecord{}, errors.New("购买渠道不能为空")
	}
	record := model.WarrantyRecord{ID: input.ID, Phone: phone, SerialNumber: serial, ExpiryDate: input.ExpiryDate, PurchaseChannel: input.PurchaseChannel, Category: input.Category, UpdatedAt: input.UpdatedAt, Active: input.Active}
	if err := s.saveRecord(ctx, record, input.ServicePoints); err != nil {
		return model.WarrantyRecord{}, err
	}
	if s.audit != nil {
		action := "create"
		if existing, lookupErr := s.reader.Lookup(ctx, phone, serial); lookupErr == nil && existing != nil {
			action = "update"
		}
		if err := s.audit.Record(ctx, model.AdminAudit{RecordID: record.ID, Action: action, Operator: operator, OccurredAt: input.UpdatedAt, Detail: fmt.Sprintf("%s %s", record.Category, record.SerialNumber)}); err != nil {
			return model.WarrantyRecord{}, err
		}
	}
	return record, nil
}

type recordSaver interface {
	SaveWarranty(context.Context, model.WarrantyRecord, []model.ServicePoint) error
}

func (s *Service) saveRecord(ctx context.Context, record model.WarrantyRecord, points []model.ServicePoint) error {
	if saver, ok := s.reader.(recordSaver); ok {
		return saver.SaveWarranty(ctx, record, points)
	}
	return errors.New("warranty reader cannot save records")
}

func (s *Service) Advice(ctx context.Context, phone, serial string) (*model.SupportAdvice, error) {
	if s.advisor == nil {
		return nil, errors.New("support advisor is unavailable")
	}
	return s.advisor.PersistMissing(ctx, model.NormalizePhone(phone), model.NormalizeSerial(serial))
}

func (s *Service) IsValidDate(value string) bool {
	return len(value) == 10 && value[4] == '-' && value[7] == '-'
}
func (s *Service) ServiceName() string { return "家电安心延保服务" }
