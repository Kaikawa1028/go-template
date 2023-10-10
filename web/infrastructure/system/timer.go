package system

import (
	"github.com/Kaikawa1028/go-template/app/domain/system"
	"time"
)

type Timer struct {
}

func NewTimer() system.ITimer {
	return &Timer{}
}

func (u Timer) Now() time.Time {
	return time.Now()
}
