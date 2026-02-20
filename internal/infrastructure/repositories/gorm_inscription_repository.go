package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormInscriptionRepository struct{ db *gorm.DB }

func NewGormInscriptionRepository(db *gorm.DB) repositories.InscriptionRepository {
	return &GormInscriptionRepository{db: db}
}

func (r *GormInscriptionRepository) FindAll(ctx context.Context, skip, take int) ([]entities.Inscription, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.InscriptionModel{}).Count(&total)

	var models []database.InscriptionModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]entities.Inscription, len(models))
	for i, m := range models {
		result[i] = toInscriptionEntity(&m)
	}
	return result, int(total), nil
}

func (r *GormInscriptionRepository) FindByID(ctx context.Context, id string) (*entities.Inscription, error) {
	var m database.InscriptionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toInscriptionEntity(&m)
	return &e, nil
}

func (r *GormInscriptionRepository) FindByUserID(ctx context.Context, userID string) ([]entities.Inscription, error) {
	var user database.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []entities.Inscription{}, nil
		}
		return nil, err
	}
	var models []database.InscriptionModel
	if err := r.db.WithContext(ctx).Where("user_ref_id = ?", user.RefID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Inscription, len(models))
	for i, m := range models {
		result[i] = toInscriptionEntity(&m)
	}
	return result, nil
}

func (r *GormInscriptionRepository) FindByTripID(ctx context.Context, tripID string) ([]entities.Inscription, error) {
	var trip database.TripModel
	if err := r.db.WithContext(ctx).Where("id = ?", tripID).First(&trip).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []entities.Inscription{}, nil
		}
		return nil, err
	}
	var models []database.InscriptionModel
	if err := r.db.WithContext(ctx).Where("trip_ref_id = ?", trip.RefID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Inscription, len(models))
	for i, m := range models {
		result[i] = toInscriptionEntity(&m)
	}
	return result, nil
}

func (r *GormInscriptionRepository) FindByIDAndUserID(ctx context.Context, id, userID string) (*entities.Inscription, error) {
	var user database.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var m database.InscriptionModel
	if err := r.db.WithContext(ctx).Where("id = ? AND user_ref_id = ?", id, user.RefID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toInscriptionEntity(&m)
	return &e, nil
}

func (r *GormInscriptionRepository) Create(ctx context.Context, data entities.CreateInscriptionData) (*entities.Inscription, error) {
	m := &database.InscriptionModel{
		UserRefID: data.UserRefID,
		TripRefID: data.TripRefID,
		Status:    "ACTIVE",
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	e := toInscriptionEntity(m)
	return &e, nil
}

func (r *GormInscriptionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.InscriptionModel{}).Error
}

func (r *GormInscriptionRepository) ExistsByUserAndTrip(ctx context.Context, userRefID, tripRefID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.InscriptionModel{}).
		Where("user_ref_id = ? AND trip_ref_id = ?", userRefID, tripRefID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormInscriptionRepository) CountByTripRefID(ctx context.Context, tripRefID int64) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.InscriptionModel{}).
		Where("trip_ref_id = ? AND status = 'ACTIVE'", tripRefID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func toInscriptionEntity(m *database.InscriptionModel) entities.Inscription {
	return entities.Inscription{
		ID: m.ID, RefID: m.RefID,
		CreatedAt: m.CreatedAt,
		UserRefID: m.UserRefID, TripRefID: m.TripRefID,
		Status: m.Status,
	}
}
