package helper

import "time"

// SetSecondZero 秒以下を0にします
func SetSecondZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}
