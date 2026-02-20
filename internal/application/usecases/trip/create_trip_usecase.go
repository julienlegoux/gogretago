package trip

import (
	"context"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateTripUseCase struct {
	tripRepository   repositories.TripRepository
	driverRepository repositories.DriverRepository
	carRepository    repositories.CarRepository
	cityRepository   repositories.CityRepository
}

func NewCreateTripUseCase(
	tripRepository repositories.TripRepository,
	driverRepository repositories.DriverRepository,
	carRepository repositories.CarRepository,
	cityRepository repositories.CityRepository,
) *CreateTripUseCase {
	return &CreateTripUseCase{
		tripRepository:   tripRepository,
		driverRepository: driverRepository,
		carRepository:    carRepository,
		cityRepository:   cityRepository,
	}
}

func (uc *CreateTripUseCase) Execute(ctx context.Context, userID string, input dtos.CreateTripInput) (*entities.Trip, error) {
	driver, err := uc.driverRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, domainerrors.NewDriverNotFoundError(userID)
	}

	car, err := uc.carRepository.FindByID(ctx, input.CarID)
	if err != nil {
		return nil, err
	}
	if car == nil {
		return nil, domainerrors.NewCarNotFoundError(input.CarID)
	}

	dateTrip, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, err
	}

	departureCity, err := uc.findOrCreateCity(ctx, input.DepartureCity)
	if err != nil {
		return nil, err
	}

	arrivalCity, err := uc.findOrCreateCity(ctx, input.ArrivalCity)
	if err != nil {
		return nil, err
	}

	return uc.tripRepository.Create(ctx, entities.CreateTripData{
		DateTrip:    dateTrip,
		Kms:         input.Kms,
		Seats:       input.Seats,
		DriverRefID: driver.RefID,
		CarRefID:    car.RefID,
		CityRefIDs:  []int64{departureCity.RefID, arrivalCity.RefID},
	})
}

func (uc *CreateTripUseCase) findOrCreateCity(ctx context.Context, cityName string) (*entities.City, error) {
	city, err := uc.cityRepository.FindByCityName(ctx, cityName)
	if err != nil {
		return nil, err
	}
	if city != nil {
		return city, nil
	}

	return uc.cityRepository.Create(ctx, entities.CreateCityData{
		CityName: cityName,
		Zipcode:  "",
	})
}
