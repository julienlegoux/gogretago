package database

import (
	"fmt"
	"time"

	"github.com/lgxju/gogretago/config"
	"github.com/lgxju/gogretago/internal/lib/shared"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AuthModel represents the authentication credentials table
type AuthModel struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID        int64      `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	Email        string     `gorm:"uniqueIndex;not null"`
	Password     string     `gorm:"not null"`
	Role         string     `gorm:"not null;default:'USER'"`
	AnonymizedAt *time.Time `gorm:"column:anonymized_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (AuthModel) TableName() string { return "auths" }

// UserModel represents the user profile table
type UserModel struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID        int64      `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	FirstName    *string    `gorm:"column:first_name"`
	LastName     *string    `gorm:"column:last_name"`
	Phone        *string    `gorm:"column:phone"`
	AuthRefID    int64      `gorm:"column:auth_ref_id;uniqueIndex;not null"`
	AnonymizedAt *time.Time `gorm:"column:anonymized_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserModel) TableName() string { return "users" }

// DriverModel represents a driver profile
type DriverModel struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID         int64      `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	DriverLicense string     `gorm:"column:driver_license;uniqueIndex;not null"`
	UserRefID     int64      `gorm:"column:user_ref_id;uniqueIndex;not null"`
	AnonymizedAt  *time.Time `gorm:"column:anonymized_at"`
}

func (DriverModel) TableName() string { return "drivers" }

// BrandModel represents a car manufacturer
type BrandModel struct {
	ID    string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID int64  `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	Name  string `gorm:"not null"`
}

func (BrandModel) TableName() string { return "brands" }

// VehicleModelModel represents a car model (e.g. "Corolla")
type VehicleModelModel struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID      int64  `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	Name       string `gorm:"not null"`
	BrandRefID int64  `gorm:"column:brand_ref_id;not null"`
}

func (VehicleModelModel) TableName() string { return "models" }

// ColorModel represents a car color
type ColorModel struct {
	ID    string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID int64  `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	Name  string `gorm:"not null"`
	Hex   string `gorm:"not null"`
}

func (ColorModel) TableName() string { return "colors" }

// ColorModelJoin represents the many-to-many between colors and models
type ColorModelJoin struct {
	ColorRefID int64 `gorm:"column:color_ref_id;primaryKey"`
	ModelRefID int64 `gorm:"column:model_ref_id;primaryKey"`
}

func (ColorModelJoin) TableName() string { return "color_models" }

// CarModel represents a car/vehicle
type CarModel struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID        int64  `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	LicensePlate string `gorm:"column:license_plate;uniqueIndex;not null"`
	ModelRefID   int64  `gorm:"column:model_ref_id;not null"`
	DriverRefID  int64  `gorm:"column:driver_ref_id;not null"`
}

func (CarModel) TableName() string { return "cars" }

// CityModel represents a city location
type CityModel struct {
	ID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID    int64  `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	CityName string `gorm:"column:city_name;not null"`
	Zipcode  string `gorm:"column:zipcode;not null;default:''"`
}

func (CityModel) TableName() string { return "cities" }

// TripModel represents a carpooling trip
type TripModel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID       int64     `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	DateTrip    time.Time `gorm:"column:date_trip;not null"`
	Kms         int       `gorm:"not null"`
	Seats       int       `gorm:"not null"`
	DriverRefID int64     `gorm:"column:driver_ref_id;not null"`
	CarRefID    int64     `gorm:"column:car_ref_id;not null"`
}

func (TripModel) TableName() string { return "trips" }

// CityTripModel represents the many-to-many between cities and trips
type CityTripModel struct {
	TripRefID int64  `gorm:"column:trip_ref_id;primaryKey"`
	CityRefID int64  `gorm:"column:city_ref_id;primaryKey"`
	Type      string `gorm:"column:type;not null"`
}

func (CityTripModel) TableName() string { return "city_trips" }

// InscriptionModel represents a passenger booking on a trip
type InscriptionModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RefID     int64     `gorm:"column:ref_id;autoIncrement;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UserRefID int64     `gorm:"column:user_ref_id;not null"`
	TripRefID int64     `gorm:"column:trip_ref_id;not null"`
	Status    string    `gorm:"not null;default:'ACTIVE'"`
}

func (InscriptionModel) TableName() string { return "inscriptions" }

var db *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() (*gorm.DB, error) {
	cfg := config.Get()
	isDevelopment := cfg.AppEnv == "development"
	appLogger := shared.NewLogger(isDevelopment)
	dbLogger := appLogger.Child(map[string]interface{}{"component": "gorm"})

	gormLogLevel := logger.Error
	if isDevelopment {
		gormLogLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: newGormLogger(dbLogger, gormLogLevel),
	}

	var err error
	db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	appLogger.Info("Database connection established", nil)
	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return db.AutoMigrate(
		&AuthModel{},
		&UserModel{},
		&DriverModel{},
		&BrandModel{},
		&VehicleModelModel{},
		&ColorModel{},
		&ColorModelJoin{},
		&CarModel{},
		&CityModel{},
		&TripModel{},
		&CityTripModel{},
		&InscriptionModel{},
	)
}
