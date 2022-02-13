package DB

import (
	_type "awesomeProject1/type"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MysqlDB *gorm.DB

func MysqlInit() {
	dsn := "root:ByteDance26!@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(fmt.Sprintf("open mysql failed, err is %s", err))
	} else {
		MysqlDB = db
	}

	//迁徙，建立表格
	if err := db.AutoMigrate(&_type.Member{}); err != nil {
		panic(err)
	}
}
