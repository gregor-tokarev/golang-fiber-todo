package models

import (
	"fmt"
	"goapi/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	if DB != nil {
		return DB
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.Cfg.PostgresHost,
		config.Cfg.PostgresUser,
		config.Cfg.PostgresPassword,
		config.Cfg.PostgresDB,
		config.Cfg.PostgresPort,
	)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}, &Task{})
	if err != nil {
		panic(err)
	}

	DB = db
	return DB
}
