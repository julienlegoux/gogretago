package inscription

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInscription_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 3}
	expectedInscription := &entities.Inscription{ID: "insc-1", RefID: 1, UserRefID: 10, TripRefID: 20}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(false, nil)
	inscriptionRepo.EXPECT().CountByTripRefID(mock.Anything, int64(20)).Return(2, nil)
	inscriptionRepo.EXPECT().Create(mock.Anything, entities.CreateInscriptionData{
		UserRefID: 10,
		TripRefID: 20,
	}).Return(expectedInscription, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInscription.ID, result.ID)
	assert.Equal(t, int64(10), result.UserRefID)
	assert.Equal(t, int64(20), result.TripRefID)
}

func TestCreateInscription_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-nonexistent"
	tripID := "trip-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(nil, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	var userNotFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &userNotFoundErr))
}

func TestCreateInscription_TripNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-nonexistent"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(nil, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	var tripNotFoundErr *domainerrors.TripNotFoundError
	assert.True(t, errors.As(err, &tripNotFoundErr))
}

func TestCreateInscription_AlreadyInscribed(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 3}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(true, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	var alreadyInscribedErr *domainerrors.AlreadyInscribedError
	assert.True(t, errors.As(err, &alreadyInscribedErr))
}

func TestCreateInscription_NoSeatsAvailable(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 3}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(false, nil)
	inscriptionRepo.EXPECT().CountByTripRefID(mock.Anything, int64(20)).Return(3, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	var noSeatsErr *domainerrors.NoSeatsAvailableError
	assert.True(t, errors.As(err, &noSeatsErr))
}

func TestCreateInscription_ExactlyAtCapacity(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 5}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(false, nil)
	// count == seats triggers the >= condition
	inscriptionRepo.EXPECT().CountByTripRefID(mock.Anything, int64(20)).Return(5, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	var noSeatsErr *domainerrors.NoSeatsAvailableError
	assert.True(t, errors.As(err, &noSeatsErr))
}

func TestCreateInscription_OneSlotLeft(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 3}
	expectedInscription := &entities.Inscription{ID: "insc-1", RefID: 1, UserRefID: 10, TripRefID: 20}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(false, nil)
	// count=2, seats=3 -> one slot left, should succeed
	inscriptionRepo.EXPECT().CountByTripRefID(mock.Anything, int64(20)).Return(2, nil)
	inscriptionRepo.EXPECT().Create(mock.Anything, entities.CreateInscriptionData{
		UserRefID: 10,
		TripRefID: 20,
	}).Return(expectedInscription, nil)

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInscription.ID, result.ID)
}

func TestCreateInscription_RepoError(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	tripID := "trip-1"

	user := &entities.PublicUser{
		User:  entities.User{ID: userID, RefID: 10},
		Email: "test@example.com",
	}
	trip := &entities.Trip{ID: tripID, RefID: 20, Seats: 3}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(user, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, tripID).Return(trip, nil)
	inscriptionRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(10), int64(20)).Return(false, nil)
	inscriptionRepo.EXPECT().CountByTripRefID(mock.Anything, int64(20)).Return(0, nil)
	inscriptionRepo.EXPECT().Create(mock.Anything, entities.CreateInscriptionData{
		UserRefID: 10,
		TripRefID: 20,
	}).Return(nil, errors.New("database error"))

	uc := NewCreateInscriptionUseCase(inscriptionRepo, userRepo, tripRepo)
	result, err := uc.Execute(ctx, userID, dtos.CreateInscriptionInput{TripID: tripID})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}
