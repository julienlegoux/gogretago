package repositories

import (
	"context"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindAll(ctx context.Context) ([]entities.PublicUser, error) {
	var users []database.UserModel
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]entities.PublicUser, 0, len(users))
	for _, u := range users {
		var auth database.AuthModel
		if err := r.db.WithContext(ctx).Where("ref_id = ?", u.AuthRefID).First(&auth).Error; err != nil {
			continue
		}
		result = append(result, toPublicUserEntity(&u, &auth))
	}
	return result, nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*entities.PublicUser, error) {
	var user database.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var auth database.AuthModel
	if err := r.db.WithContext(ctx).Where("ref_id = ?", user.AuthRefID).First(&auth).Error; err != nil {
		return nil, err
	}

	pu := toPublicUserEntity(&user, &auth)
	return &pu, nil
}

func (r *GormUserRepository) FindByAuthRefID(ctx context.Context, authRefID int64) (*entities.PublicUser, error) {
	var user database.UserModel
	if err := r.db.WithContext(ctx).Where("auth_ref_id = ?", authRefID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var auth database.AuthModel
	if err := r.db.WithContext(ctx).Where("ref_id = ?", authRefID).First(&auth).Error; err != nil {
		return nil, err
	}

	pu := toPublicUserEntity(&user, &auth)
	return &pu, nil
}

func (r *GormUserRepository) Update(ctx context.Context, id string, data entities.UpdateUserData) (*entities.PublicUser, error) {
	updates := map[string]interface{}{}
	if data.FirstName != nil {
		updates["first_name"] = *data.FirstName
	}
	if data.LastName != nil {
		updates["last_name"] = *data.LastName
	}
	if data.Phone != nil {
		updates["phone"] = *data.Phone
	}

	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&database.UserModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return r.FindByID(ctx, id)
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.UserModel{}).Error
}

func (r *GormUserRepository) Anonymize(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&database.UserModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"first_name":    nil,
		"last_name":     nil,
		"phone":         nil,
		"anonymized_at": now,
	}).Error
}

func toPublicUserEntity(model *database.UserModel, auth *database.AuthModel) entities.PublicUser {
	return entities.PublicUser{
		User: entities.User{
			ID:           model.ID,
			RefID:        model.RefID,
			FirstName:    model.FirstName,
			LastName:     model.LastName,
			Phone:        model.Phone,
			AuthRefID:    model.AuthRefID,
			AnonymizedAt: model.AnonymizedAt,
			CreatedAt:    model.CreatedAt,
			UpdatedAt:    model.UpdatedAt,
		},
		Email: auth.Email,
	}
}
