//go:build integration

package handler_test

import (
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/Kaikawa1028/go-template/app/domain/system"
	"github.com/Kaikawa1028/go-template/app/infrastructure/persistence"
	"github.com/Kaikawa1028/go-template/app/wire"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Kaikawa1028/go-template/app/config"
)

var loggerHook *test.Hook

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", "..", ".env.integration"))

	loggerHook = test.NewGlobal()

	dbConfig, err := config.NewDBConfig()
	if err != nil {
		panic(err)
	}

	persistence.RegisterTxdbDriver(dbConfig)

	code := m.Run()
	os.Exit(code)
}

type TestCaseContainer struct {
	t         *testing.T
	DI        *wire.IntegrationTestDIContainer
	Cleanup   func()
	MockCtrl  *gomock.Controller
	MockTimer *system.MockITimer
}

// テストケース毎に実行される前処理
func beforeEach(t *testing.T) *TestCaseContainer {
	ctrl := gomock.NewController(t)

	loggerHook.Reset()

	mockTimer := system.NewMockITimer(ctrl)

	di, cleanup, err := wire.InitializeIntegrationTestDIContainer(mockTimer)
	if err != nil {
		t.Fatal(err)
	}

	return &TestCaseContainer{
		t:         t,
		DI:        di,
		Cleanup:   cleanup,
		MockCtrl:  ctrl,
		MockTimer: mockTimer,
	}
}

// テストケース毎に実行される後処理
func afterEach(cont *TestCaseContainer) {
	cont.Cleanup()
	cont.MockCtrl.Finish()
}

// prepareTestData 外部キー制約のチェックを無効化した状態で第二引数の処理を実行します
func prepareTestData(db *gorm.DB, closure func(db *gorm.DB)) {
	mustExec(db, "SET FOREIGN_KEY_CHECKS = 0")
	closure(db)
	mustExec(db, "SET FOREIGN_KEY_CHECKS = 1")
}

// mustInsert データを挿入し、エラーが発生した場合はpanicを発生させます
func mustInsert(db *gorm.DB, table string, records []map[string]interface{}) {
	err := db.Table(table).Create(records).Error
	if err != nil {
		panic(err)
	}
}

// mustExec SQLを実行し、エラーが発生した場合はpanicを発生させます
func mustExec(db *gorm.DB, sql string) {
	err := db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}

// mustNewDateTime 日時を生成します
func mustNewDateTime(dateTime string) time.Time {
	now, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		panic(err)
	}
	return now
}

// mustOpenFile 日時を生成します
func mustOpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

// assertFileSize ファイルサイズを確認します（sizeの単位：バイト）
func assertFileSize(t *testing.T, path string, size int64) {
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, size, fi.Size())
}
