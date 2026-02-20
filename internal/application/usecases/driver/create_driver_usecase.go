package driver

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateDriverUseCase struct {
	driverRepository repositories.DriverRepository
	userRepository   repositories.UserRepository
	authRepository   repositories.AuthRepository
}

func NewCreateDriverUseCase(
	driverRepository repositories.DriverRepository,
	userRepository repositories.UserRepository,
	authRepository repositories.AuthRepository,
) *CreateDriverUseCase {
	return &CreateDriverUseCase{
		driverRepository: driverRepository,
		userRepository:   userRepository,
		authRepository:   authRepository,
	}
}

func (uc *CreateDriverUseCase) Execute(ctx context.Context, userID string, input dtos.CreateDriverInput) (*entities.Driver, error) {
	user, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewUserNotFoundError(userID)
	}

	existingDriver, err := uc.driverRepository.FindByUserRefID(ctx, user.RefID)
	if err != nil {
		return nil, err
	}
	if existingDriver != nil {
		return nil, domainerrors.NewDriverAlreadyExistsError(userID)
	}

	driver, err := uc.driverRepository.Create(ctx, entities.CreateDriverData{
		DriverLicense: input.DriverLicense,
		UserRefID:     user.RefID,
	})
	if err != nil {
		return nil, err
	}

	err = uc.authRepository.UpdateRole(ctx, user.AuthRefID, "DRIVER")
	if err != nil {
		return nil, err
	}

	return driver, nil
}
