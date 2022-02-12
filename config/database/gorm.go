package database

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"perkawis/src/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type gormInstance struct {
	master *gorm.DB
}

// Master initialize DB for master data
func (g *gormInstance) Master() *gorm.DB {
	return g.master
}

// GormDatabase abstraction
type GormDatabase interface {
	Master() *gorm.DB
}

func InitGorm() GormDatabase {
	inst := new(gormInstance)

	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	// username, password, host, port, database
	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DBNAME"))

	fmt.Println(connection)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: connection,
	}), gormConfig)

	if err != nil {
		fmt.Printf("cant connect to database: %s", err)
		panic("connection error")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Infof("failed to load generic database object")
	}

	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetMaxIdleConns(10000)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	db.AutoMigrate(&model.Activity{}, &model.Todo{})

	inst.master = db

	return inst
}
