package user

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type AnonymizeUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewAnonymizeUserUseCase(userRepository repositories.UserRepository) *AnonymizeUserUseCase {
	return &AnonymizeUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *AnonymizeUserUseCase) Execute(ctx context.Context, id string) error {
	user, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return domainerrors.NewUserNotFoundError(id)
	}

	return uc.userRepository.Anonymize(ctx, id)
}
