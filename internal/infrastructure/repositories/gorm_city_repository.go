package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormCityRepository struct{ db *gorm.DB }

func NewGormCityRepository(db *gorm.DB) repositories.CityRepository {
	return &GormCityRepository{db: db}
}

func (r *GormCityRepository) FindAll(ctx context.Context, skip, take int) ([]entities.City, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.CityModel{}).Count(&total)

	var models []database.CityModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]entities.City, len(models))
	for i, m := range models {
		result[i] = toCityEntity(&m)
	}
	return result, int(total), nil
}

func (r *GormCityRepository) FindByID(ctx context.Context, id string) (*entities.City, error) {
	var m database.CityModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toCityEntity(&m)
	return &e, nil
}

func (r *GormCityRepository) FindByCityName(ctx context.Context, name string) (*entities.City, error) {
	var m database.CityModel
	if err := r.db.WithContext(ctx).Where("city_name = ?", name).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toCityEntity(&m)
	return &e, nil
}

func (r *GormCityRepository) Create(ctx context.Context, data entities.CreateCityData) (*entities.City, error) {
	m := &database.CityModel{CityName: data.CityName, Zipcode: data.Zipcode}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	e := toCityEntity(m)
	return &e, nil
}

func (r *GormCityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.CityModel{}).Error
}

func toCityEntity(m *database.CityModel) entities.City {
	return entities.City{ID: m.ID, RefID: m.RefID, CityName: m.CityName, Zipcode: m.Zipcode}
}
