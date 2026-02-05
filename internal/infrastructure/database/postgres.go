package database

import (
	"fmt"
	"log"

	"github.com/lgxju/gogretago/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// UserModel represents the database model for users
type UserModel struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"column:first_name;not null"`
	LastName  string `gorm:"column:last_name;not null"`
	Phone     string `gorm:"not null"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

var db *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() (*gorm.DB, error) {
	cfg := config.Get()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}

	if cfg.AppEnv == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var err error
	db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return db.AutoMigrate(&UserModel{})
}
