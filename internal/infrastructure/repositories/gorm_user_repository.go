package repositories

import (
	"context"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &GormUserRepository{db: db}
}

// FindByID finds a user by ID
func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var model database.UserModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toEntity(&model), nil
}

// FindByEmail finds a user by email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model database.UserModel
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toEntity(&model), nil
}

// Create creates a new user
func (r *GormUserRepository) Create(ctx context.Context, data entities.CreateUserData) (*entities.User, error) {
	model := &database.UserModel{
		Email:     data.Email,
		Password:  data.Password,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Phone:     data.Phone,
	}

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return toEntity(model), nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *GormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&database.UserModel{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// toEntity converts a database model to a domain entity
func toEntity(model *database.UserModel) *entities.User {
	return &entities.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Phone:     model.Phone,
		CreatedAt: time.Unix(model.CreatedAt, 0),
		UpdatedAt: time.Unix(model.UpdatedAt, 0),
	}
}
