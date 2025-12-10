package db

/*
初始化GORM 连接Postgres AutoMigrate
*/

import(
	"fmt"
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB 初始化并返回 *gorm.DB
func NewPostgresDB() *gorm.DB {
	// 先用硬编码，后面再从env中读取
	host := "localhost"
	port := 5432
	user := "root"
	password := "root"
	dbname := "go_tasker"

	// TimeZone is the correct DSN parameter (was misspelled as TimeZeno).
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "[gorm] ", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful: true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&TaskModel{}, &UserModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
