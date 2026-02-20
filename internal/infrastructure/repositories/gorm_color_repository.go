package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormColorRepository struct{ db *gorm.DB }

func NewGormColorRepository(db *gorm.DB) repositories.ColorRepository {
	return &GormColorRepository{db: db}
}

func (r *GormColorRepository) FindAll(ctx context.Context, skip, take int) ([]entities.Color, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.ColorModel{}).Count(&total)

	var models []database.ColorModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]entities.Color, len(models))
	for i, m := range models {
		result[i] = entities.Color{ID: m.ID, RefID: m.RefID, Name: m.Name, Hex: m.Hex}
	}
	return result, int(total), nil
}

func (r *GormColorRepository) FindByID(ctx context.Context, id string) (*entities.Color, error) {
	var m database.ColorModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.Color{ID: m.ID, RefID: m.RefID, Name: m.Name, Hex: m.Hex}, nil
}

func (r *GormColorRepository) FindByName(ctx context.Context, name string) (*entities.Color, error) {
	var m database.ColorModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.Color{ID: m.ID, RefID: m.RefID, Name: m.Name, Hex: m.Hex}, nil
}

func (r *GormColorRepository) Create(ctx context.Context, data entities.CreateColorData) (*entities.Color, error) {
	m := &database.ColorModel{Name: data.Name, Hex: data.Hex}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return &entities.Color{ID: m.ID, RefID: m.RefID, Name: m.Name, Hex: m.Hex}, nil
}

func (r *GormColorRepository) Update(ctx context.Context, id string, data entities.UpdateColorData) (*entities.Color, error) {
	updates := map[string]interface{}{}
	if data.Name != nil {
		updates["name"] = *data.Name
	}
	if data.Hex != nil {
		updates["hex"] = *data.Hex
	}
	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&database.ColorModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}

func (r *GormColorRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.ColorModel{}).Error
}
