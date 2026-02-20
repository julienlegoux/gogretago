package user

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type GetUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetUserUseCase(userRepository repositories.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*entities.PublicUser, error) {
	user, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewUserNotFoundError(id)
	}
	if user.AnonymizedAt != nil {
		return nil, domainerrors.NewUserNotFoundError(id)
	}
	return user, nil
}
