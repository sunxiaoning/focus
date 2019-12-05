package app

import (
	"fmt"
	"focus/cfg"
	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

// 初始化数据库
func InitDB() error {
	if err := configDB(); err != nil {
		return err
	}
	return nil
}

func ReConnectDB() {
	err := backoff.Retry(reConnectDB, backoff.NewConstantBackOff(time.Second*2))
	if err != nil {
		logrus.Error("reConnectDB error!", err)
	}
}

func reConnectDB() error {
	logrus.Info("start to check DB status...")
	var err error = nil
	if err = cfg.FocusCtx.DB.DB().Ping(); err != nil {
		err = configDB()
	}
	logrus.Info("check DB status end...")
	return err
}

func configDB() error {
	oldCfg := cfg.FocusCtx.Cfg
	dbConfig := oldCfg.Database
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName))
	if err != nil {
		return err
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	testDB(db)
	cfg.FocusCtx.DB = db
	return nil
}

func testDB(db *gorm.DB) {
	type Customer struct {
		ID          int64
		EnName      string
		ChineseName string
	}
	var customer Customer
	db.First(&customer, 999)
	logrus.Printf("%s=%s", customer.EnName, customer.ChineseName)
}
