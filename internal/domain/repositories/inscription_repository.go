package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// InscriptionRepository defines the interface for inscription persistence operations
type InscriptionRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.Inscription, int, error)
	FindByID(ctx context.Context, id string) (*entities.Inscription, error)
	FindByUserID(ctx context.Context, userID string) ([]entities.Inscription, error)
	FindByTripID(ctx context.Context, tripID string) ([]entities.Inscription, error)
	FindByIDAndUserID(ctx context.Context, id string, userID string) (*entities.Inscription, error)
	Create(ctx context.Context, data entities.CreateInscriptionData) (*entities.Inscription, error)
	Delete(ctx context.Context, id string) error
	ExistsByUserAndTrip(ctx context.Context, userRefID, tripRefID int64) (bool, error)
	CountByTripRefID(ctx context.Context, tripRefID int64) (int, error)
}
