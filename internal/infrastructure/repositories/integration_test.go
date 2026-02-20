//go:build integration

package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate postgres container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	testDB, err = gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Enable gen_random_uuid()
	testDB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)

	// Run migrations
	if err := testDB.AutoMigrate(
		&database.AuthModel{},
		&database.UserModel{},
		&database.DriverModel{},
		&database.BrandModel{},
		&database.VehicleModelModel{},
		&database.ColorModel{},
		&database.ColorModelJoin{},
		&database.CarModel{},
		&database.CityModel{},
		&database.TripModel{},
		&database.CityTripModel{},
		&database.InscriptionModel{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Test database ready")
	os.Exit(m.Run())
}

func cleanTables(t *testing.T) {
	t.Helper()
	tables := []string{
		"inscriptions",
		"city_trips",
		"trips",
		"cars",
		"color_models",
		"colors",
		"models",
		"brands",
		"cities",
		"drivers",
		"users",
		"auths",
	}
	for _, table := range tables {
		if err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
