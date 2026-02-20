package user

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListUsersUseCase struct {
	userRepository repositories.UserRepository
}

func NewListUsersUseCase(userRepository repositories.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepository: userRepository,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]entities.PublicUser, error) {
	return uc.userRepository.FindAll(ctx)
}
