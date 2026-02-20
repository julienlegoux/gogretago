package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// CarRepository defines the interface for car persistence operations
type CarRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.Car, int, error)
	FindByID(ctx context.Context, id string) (*entities.Car, error)
	Create(ctx context.Context, data entities.CreateCarData) (*entities.Car, error)
	Update(ctx context.Context, id string, data entities.UpdateCarData) (*entities.Car, error)
	Delete(ctx context.Context, id string) error
	ExistsByLicensePlate(ctx context.Context, licensePlate string) (bool, error)
}
