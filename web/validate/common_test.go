package validate_test

import (
	"github.com/golang/mock/gomock"
	"github.com/Kaikawa1028/go-template/app/domain/system"
	"testing"
	"time"
)

type TestCaseContainer struct {
	t         *testing.T
	MockCtrl  *gomock.Controller
	MockTimer *system.MockITimer
}

// beforeEach テストケース毎に実行される前処理
func beforeEach(t *testing.T) *TestCaseContainer {
	ctrl := gomock.NewController(t)

	mockTimer := system.NewMockITimer(ctrl)

	return &TestCaseContainer{
		t:         t,
		MockCtrl:  ctrl,
		MockTimer: mockTimer,
	}
}

// afterEach テストケース毎に実行される後処理
func afterEach(cont *TestCaseContainer) {
	cont.MockCtrl.Finish()
}

// mustNewDateTime 日時を生成します
func mustNewDateTime(dateTime string) time.Time {
	now, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		panic(err)
	}
	return now
}
