package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormAuthRepository struct {
	db *gorm.DB
}

func NewGormAuthRepository(db *gorm.DB) repositories.AuthRepository {
	return &GormAuthRepository{db: db}
}

func (r *GormAuthRepository) FindByEmail(ctx context.Context, email string) (*entities.Auth, error) {
	var model database.AuthModel
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toAuthEntity(&model), nil
}

func (r *GormAuthRepository) CreateWithUser(ctx context.Context, authData entities.CreateAuthData, userData entities.CreateUserData) (*entities.Auth, *entities.PublicUser, error) {
	var auth *entities.Auth
	var publicUser *entities.PublicUser

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		authModel := &database.AuthModel{
			Email:    authData.Email,
			Password: authData.Password,
			Role:     "USER",
		}
		if err := tx.Create(authModel).Error; err != nil {
			return err
		}

		userModel := &database.UserModel{
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Phone:     userData.Phone,
			AuthRefID: authModel.RefID,
		}
		if err := tx.Create(userModel).Error; err != nil {
			return err
		}

		auth = toAuthEntity(authModel)
		publicUser = &entities.PublicUser{
			User: entities.User{
				ID:        userModel.ID,
				RefID:     userModel.RefID,
				FirstName: userModel.FirstName,
				LastName:  userModel.LastName,
				Phone:     userModel.Phone,
				AuthRefID: userModel.AuthRefID,
				CreatedAt: userModel.CreatedAt,
				UpdatedAt: userModel.UpdatedAt,
			},
			Email: authModel.Email,
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return auth, publicUser, nil
}

func (r *GormAuthRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&database.AuthModel{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *GormAuthRepository) UpdateRole(ctx context.Context, refID int64, role string) error {
	return r.db.WithContext(ctx).Model(&database.AuthModel{}).Where("ref_id = ?", refID).Update("role", role).Error
}

func toAuthEntity(model *database.AuthModel) *entities.Auth {
	return &entities.Auth{
		ID:           model.ID,
		RefID:        model.RefID,
		Email:        model.Email,
		Password:     model.Password,
		Role:         model.Role,
		AnonymizedAt: model.AnonymizedAt,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}
