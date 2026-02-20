package trip

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTrip_Success(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}
	car := &entities.Car{
		ID:           "car-1",
		RefID:        400,
		LicensePlate: "ABC-123",
	}
	departureCity := &entities.City{
		ID:       "city-1",
		RefID:    10,
		CityName: "Paris",
		Zipcode:  "75000",
	}
	arrivalCity := &entities.City{
		ID:       "city-2",
		RefID:    20,
		CityName: "Lyon",
		Zipcode:  "69000",
	}

	dateTrip, _ := time.Parse("2006-01-02", "2026-06-15")
	createdTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		DateTrip:    dateTrip,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
	}

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	carRepo.EXPECT().FindByID(ctx, "car-1").Return(car, nil)
	cityRepo.EXPECT().FindByCityName(ctx, "Paris").Return(departureCity, nil)
	cityRepo.EXPECT().FindByCityName(ctx, "Lyon").Return(arrivalCity, nil)
	tripRepo.EXPECT().Create(ctx, entities.CreateTripData{
		DateTrip:    dateTrip,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
		CityRefIDs:  []int64{10, 20},
	}).Return(createdTrip, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "2026-06-15",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "trip-1", result.ID)
	assert.Equal(t, 450, result.Kms)
	assert.Equal(t, 3, result.Seats)
	assert.Equal(t, int64(300), result.DriverRefID)
}

func TestCreateTrip_DriverNotFound(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(nil, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "2026-06-15",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.DriverNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestCreateTrip_CarNotFound(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	carRepo.EXPECT().FindByID(ctx, "car-1").Return(nil, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "2026-06-15",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.CarNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestCreateTrip_InvalidDateFormat(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}
	car := &entities.Car{
		ID:           "car-1",
		RefID:        400,
		LicensePlate: "ABC-123",
	}

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	carRepo.EXPECT().FindByID(ctx, "car-1").Return(car, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "not-a-date",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestCreateTrip_FindsExistingCity(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}
	car := &entities.Car{
		ID:           "car-1",
		RefID:        400,
		LicensePlate: "ABC-123",
	}
	existingDeparture := &entities.City{
		ID:       "city-1",
		RefID:    10,
		CityName: "Paris",
		Zipcode:  "75000",
	}
	existingArrival := &entities.City{
		ID:       "city-2",
		RefID:    20,
		CityName: "Lyon",
		Zipcode:  "69000",
	}

	dateTrip, _ := time.Parse("2006-01-02", "2026-06-15")
	createdTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		DateTrip:    dateTrip,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
	}

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	carRepo.EXPECT().FindByID(ctx, "car-1").Return(car, nil)
	// Both cities already exist - no Create calls expected
	cityRepo.EXPECT().FindByCityName(ctx, "Paris").Return(existingDeparture, nil)
	cityRepo.EXPECT().FindByCityName(ctx, "Lyon").Return(existingArrival, nil)
	tripRepo.EXPECT().Create(ctx, entities.CreateTripData{
		DateTrip:    dateTrip,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
		CityRefIDs:  []int64{10, 20},
	}).Return(createdTrip, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "2026-06-15",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "trip-1", result.ID)
}

func TestCreateTrip_CreatesNewCity(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}
	car := &entities.Car{
		ID:           "car-1",
		RefID:        400,
		LicensePlate: "ABC-123",
	}
	newDeparture := &entities.City{
		ID:       "city-new-1",
		RefID:    30,
		CityName: "Marseille",
		Zipcode:  "",
	}
	existingArrival := &entities.City{
		ID:       "city-2",
		RefID:    20,
		CityName: "Lyon",
		Zipcode:  "69000",
	}

	dateTrip, _ := time.Parse("2006-01-02", "2026-06-15")
	createdTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		DateTrip:    dateTrip,
		Kms:         350,
		Seats:       2,
		DriverRefID: 300,
		CarRefID:    400,
	}

	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	carRepo.EXPECT().FindByID(ctx, "car-1").Return(car, nil)
	// Departure city not found, so it gets created with empty zipcode
	cityRepo.EXPECT().FindByCityName(ctx, "Marseille").Return(nil, nil)
	cityRepo.EXPECT().Create(ctx, entities.CreateCityData{
		CityName: "Marseille",
		Zipcode:  "",
	}).Return(newDeparture, nil)
	cityRepo.EXPECT().FindByCityName(ctx, "Lyon").Return(existingArrival, nil)
	tripRepo.EXPECT().Create(ctx, entities.CreateTripData{
		DateTrip:    dateTrip,
		Kms:         350,
		Seats:       2,
		DriverRefID: 300,
		CarRefID:    400,
		CityRefIDs:  []int64{30, 20},
	}).Return(createdTrip, nil)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           350,
		Date:          "2026-06-15",
		DepartureCity: "Marseille",
		ArrivalCity:   "Lyon",
		Seats:         2,
		CarID:         "car-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "trip-1", result.ID)
	assert.Equal(t, 350, result.Kms)
}

func TestCreateTrip_RepoError(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	repoErr := errors.New("database error")
	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(nil, repoErr)

	uc := NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateTripInput{
		Kms:           450,
		Date:          "2026-06-15",
		DepartureCity: "Paris",
		ArrivalCity:   "Lyon",
		Seats:         3,
		CarID:         "car-1",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

// Ensure mock import is used
var _ mock.TestingT = (*testing.T)(nil)
