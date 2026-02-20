package user

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type UpdateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewUpdateUserUseCase(userRepository repositories.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, userID string, input dtos.UpdateProfileInput) (*entities.PublicUser, error) {
	existing, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domainerrors.NewUserNotFoundError(userID)
	}

	updated, err := uc.userRepository.Update(ctx, userID, entities.UpdateUserData{
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
		Phone:     &input.Phone,
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}
