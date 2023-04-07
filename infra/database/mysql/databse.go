package mysql

import (
	"fmt"

	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func buildDSN(user, password, dbName, address, port string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", user, password, address, port, dbName)
}

func NewDatabase(conf *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(conf.DBConf.UserName, conf.DBConf.Password, conf.DBConf.DBName, conf.DBConf.Addr, conf.DBConf.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
