package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormDriverRepository struct{ db *gorm.DB }

func NewGormDriverRepository(db *gorm.DB) repositories.DriverRepository {
	return &GormDriverRepository{db: db}
}

func (r *GormDriverRepository) FindByUserRefID(ctx context.Context, userRefID int64) (*entities.Driver, error) {
	var m database.DriverModel
	if err := r.db.WithContext(ctx).Where("user_ref_id = ?", userRefID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toDriverEntity(&m), nil
}

func (r *GormDriverRepository) FindByUserID(ctx context.Context, userID string) (*entities.Driver, error) {
	var user database.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.FindByUserRefID(ctx, user.RefID)
}

func (r *GormDriverRepository) Create(ctx context.Context, data entities.CreateDriverData) (*entities.Driver, error) {
	m := &database.DriverModel{
		DriverLicense: data.DriverLicense,
		UserRefID:     data.UserRefID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return toDriverEntity(m), nil
}

func toDriverEntity(m *database.DriverModel) *entities.Driver {
	return &entities.Driver{
		ID:            m.ID,
		RefID:         m.RefID,
		DriverLicense: m.DriverLicense,
		UserRefID:     m.UserRefID,
		AnonymizedAt:  m.AnonymizedAt,
	}
}
