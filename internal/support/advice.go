package support

import (
	"context"
	"errors"
	"fmt"

	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
)

type Advisor struct {
	store     *persistence.Store
	hotline   string
	channels  string
	fixedDate string
}

func NewAdvisor(store *persistence.Store, hotline, channels, fixedDate string) *Advisor {
	return &Advisor{store: store, hotline: hotline, channels: channels, fixedDate: fixedDate}
}

func (a *Advisor) BuildMissingAdvice(phone, serial string) *model.SupportAdvice {
	return &model.SupportAdvice{Phone: phone, SerialNumber: serial, Message: "暂未找到该设备的延保记录，请核对购买手机号和设备序列号后重试；如仍有疑问请联系客服。", CreatedAt: a.fixedDate, Hotline: a.hotline, Channels: a.channels}
}

func (a *Advisor) PersistMissing(ctx context.Context, phone, serial string) (*model.SupportAdvice, error) {
	if a == nil {
		return nil, errors.New("advisor is unavailable")
	}
	advice := a.BuildMissingAdvice(phone, serial)
	if a.store != nil {
		if err := a.store.SaveAdvice(ctx, *advice); err != nil {
			return nil, fmt.Errorf("save support advice: %w", err)
		}
	}
	return advice, nil
}

func (a *Advisor) GetLatest(ctx context.Context, phone, serial string) (*model.SupportAdvice, error) {
	if a == nil || a.store == nil {
		return nil, errors.New("advisor storage is unavailable")
	}
	return a.store.LatestAdvice(ctx, phone, serial)
}

func (a *Advisor) ValidateAdvice(advice *model.SupportAdvice) error {
	if advice == nil {
		return errors.New("advice is nil")
	}
	if advice.Phone == "" || advice.SerialNumber == "" {
		return errors.New("advice identity is incomplete")
	}
	if advice.Message == "" {
		return errors.New("advice message is empty")
	}
	if advice.Hotline == "" {
		return errors.New("advice hotline is empty")
	}
	return nil
}
