package persistence

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/Kaikawa1028/go-template/app/config"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/logger"
)

// NewOrmConnection gormによるDB接続を開始します
func NewOrmConnection(config *config.DBConfig) (*gorm.DB, func(), error) {
	loc := strings.Replace(config.TimeZone, "/", "%2F", 1)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s", config.User, config.Password, config.Host, config.Port, config.Database, loc)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent), // gormによるログ出力を抑制する
	})
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, nil, errors.Wrap(err)
	}

	cleanup := func() {
		sqlDB.Close()
	}

	logger.SimpleInfoF("Connected to the database with gorm %s:%d/%s...", config.Host, config.Port, config.Database)
	return db, cleanup, nil
}

// RegisterTxdbDriver 自動でロールバックする単一トランザクションのDBドライバを登録する（テスト用）
func RegisterTxdbDriver(config *config.DBConfig) {
	loc := strings.Replace(config.TimeZone, "/", "%2F", 1)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s", config.User, config.Password, config.Host, config.Port, config.Database, loc)
	txdb.Register("txdb", "mysql", dsn)
}

// NewOrmConnectionWithTxdb テスト用のドライバでDBに接続する
func NewOrmConnectionWithTxdb() (*gorm.DB, func(), error) {
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DriverName: "txdb",
		}),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent), // gormによるログ出力を抑制する
		},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}
	cleanup := func() {
		sqlDB.Close()
	}
	return db, cleanup, nil
}
