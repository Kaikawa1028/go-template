package logger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/Kaikawa1028/go-template/app/logger"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var loggerHook *test.Hook

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", ".env.integration"))

	err := logger.SetupLogger()
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(io.Discard)

	loggerHook = test.NewGlobal()

	code := m.Run()
	os.Exit(code)
}

type TestCaseContainer struct {
	t        *testing.T
	MockCtrl *gomock.Controller
}

// テストケース毎に実行される前処理
func beforeEach(t *testing.T) *TestCaseContainer {
	ctrl := gomock.NewController(t)

	loggerHook.Reset()

	return &TestCaseContainer{
		t:        t,
		MockCtrl: ctrl,
	}
}

// テストケース毎に実行される後処理
func afterEach(cont *TestCaseContainer) {
	cont.MockCtrl.Finish()
}
