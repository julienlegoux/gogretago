package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormBrandRepository struct{ db *gorm.DB }

func NewGormBrandRepository(db *gorm.DB) repositories.BrandRepository {
	return &GormBrandRepository{db: db}
}

func (r *GormBrandRepository) FindAll(ctx context.Context, skip, take int) ([]entities.Brand, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.BrandModel{}).Count(&total)

	var models []database.BrandModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.Brand, len(models))
	for i, m := range models {
		result[i] = entities.Brand{ID: m.ID, RefID: m.RefID, Name: m.Name}
	}
	return result, int(total), nil
}

func (r *GormBrandRepository) FindByID(ctx context.Context, id string) (*entities.Brand, error) {
	var m database.BrandModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.Brand{ID: m.ID, RefID: m.RefID, Name: m.Name}, nil
}

func (r *GormBrandRepository) Create(ctx context.Context, data entities.CreateBrandData) (*entities.Brand, error) {
	m := &database.BrandModel{Name: data.Name}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return &entities.Brand{ID: m.ID, RefID: m.RefID, Name: m.Name}, nil
}

func (r *GormBrandRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.BrandModel{}).Error
}
