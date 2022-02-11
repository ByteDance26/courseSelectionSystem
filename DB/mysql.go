package DB

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MysqlDB *gorm.DB

func MysqlInit() {
	dsn := "root:ByteDance26!@tcp(180.184.70.231:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(fmt.Sprintf("open mysql failed, err is %s", err))
	} else {
		MysqlDB = db
	}
}
