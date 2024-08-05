package domain

import (
	"context"
	"go.uber.org/zap"
)

const (
	SendedNotifyStatus   = "send"
	NoSendedNotifyStatus = "no send"
)

type Notify struct {
	ID       int
	FlatID   int
	HouseID  int
	UserMail string
	Status   string
}

type NotifySender interface {
	SendEmail(ctx context.Context, recipient string, message string) error
}

type NotifyRepo interface {
	GetNoSendNotifies(ctx context.Context, lg *zap.Logger) ([]Notify, error)
	SendNotifyByID(ctx context.Context, id int, lg *zap.Logger) error
}
