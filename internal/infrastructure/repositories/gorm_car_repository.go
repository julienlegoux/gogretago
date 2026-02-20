package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormCarRepository struct{ db *gorm.DB }

func NewGormCarRepository(db *gorm.DB) repositories.CarRepository {
	return &GormCarRepository{db: db}
}

func (r *GormCarRepository) FindAll(ctx context.Context, skip, take int) ([]entities.Car, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.CarModel{}).Count(&total)

	var models []database.CarModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]entities.Car, len(models))
	for i, m := range models {
		result[i] = toCarEntity(&m)
	}
	return result, int(total), nil
}

func (r *GormCarRepository) FindByID(ctx context.Context, id string) (*entities.Car, error) {
	var m database.CarModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toCarEntity(&m)
	return &e, nil
}

func (r *GormCarRepository) Create(ctx context.Context, data entities.CreateCarData) (*entities.Car, error) {
	m := &database.CarModel{
		LicensePlate: data.LicensePlate,
		ModelRefID:   data.ModelRefID,
		DriverRefID:  data.DriverRefID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	e := toCarEntity(m)
	return &e, nil
}

func (r *GormCarRepository) Update(ctx context.Context, id string, data entities.UpdateCarData) (*entities.Car, error) {
	updates := map[string]interface{}{}
	if data.LicensePlate != nil {
		updates["license_plate"] = *data.LicensePlate
	}
	if data.ModelRefID != nil {
		updates["model_ref_id"] = *data.ModelRefID
	}
	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&database.CarModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}

func (r *GormCarRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.CarModel{}).Error
}

func (r *GormCarRepository) ExistsByLicensePlate(ctx context.Context, licensePlate string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.CarModel{}).Where("license_plate = ?", licensePlate).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toCarEntity(m *database.CarModel) entities.Car {
	return entities.Car{
		ID: m.ID, RefID: m.RefID,
		LicensePlate: m.LicensePlate,
		ModelRefID:   m.ModelRefID,
		DriverRefID:  m.DriverRefID,
	}
}
