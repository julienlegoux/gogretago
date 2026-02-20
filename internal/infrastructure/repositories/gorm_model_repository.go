package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormModelRepository struct{ db *gorm.DB }

func NewGormModelRepository(db *gorm.DB) repositories.ModelRepository {
	return &GormModelRepository{db: db}
}

func (r *GormModelRepository) FindAll(ctx context.Context) ([]entities.VehicleModel, error) {
	var models []database.VehicleModelModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.VehicleModel, len(models))
	for i, m := range models {
		result[i] = toVehicleModelEntity(&m)
	}
	return result, nil
}

func (r *GormModelRepository) FindByID(ctx context.Context, id string) (*entities.VehicleModel, error) {
	var m database.VehicleModelModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toVehicleModelEntity(&m)
	return &e, nil
}

func (r *GormModelRepository) FindByNameAndBrand(ctx context.Context, name string, brandRefID int64) (*entities.VehicleModel, error) {
	var m database.VehicleModelModel
	if err := r.db.WithContext(ctx).Where("name = ? AND brand_ref_id = ?", name, brandRefID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toVehicleModelEntity(&m)
	return &e, nil
}

func (r *GormModelRepository) Create(ctx context.Context, data entities.CreateModelData) (*entities.VehicleModel, error) {
	m := &database.VehicleModelModel{Name: data.Name, BrandRefID: data.BrandRefID}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	e := toVehicleModelEntity(m)
	return &e, nil
}

func toVehicleModelEntity(m *database.VehicleModelModel) entities.VehicleModel {
	return entities.VehicleModel{ID: m.ID, RefID: m.RefID, Name: m.Name, BrandRefID: m.BrandRefID}
}
