package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GormTripRepository struct{ db *gorm.DB }

func NewGormTripRepository(db *gorm.DB) repositories.TripRepository {
	return &GormTripRepository{db: db}
}

func (r *GormTripRepository) FindAll(ctx context.Context, skip, take int) ([]entities.Trip, int, error) {
	var total int64
	r.db.WithContext(ctx).Model(&database.TripModel{}).Count(&total)

	var models []database.TripModel
	if err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]entities.Trip, len(models))
	for i, m := range models {
		result[i] = toTripEntity(&m)
	}
	return result, int(total), nil
}

func (r *GormTripRepository) FindByID(ctx context.Context, id string) (*entities.Trip, error) {
	var m database.TripModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	e := toTripEntity(&m)
	return &e, nil
}

func (r *GormTripRepository) FindByFilters(ctx context.Context, filters entities.TripFilters) ([]entities.Trip, error) {
	query := r.db.WithContext(ctx).Model(&database.TripModel{})

	if filters.DepartureCity != nil {
		query = query.Where("ref_id IN (SELECT trip_ref_id FROM city_trips ct JOIN cities c ON c.ref_id = ct.city_ref_id WHERE ct.type = 'DEPARTURE' AND c.city_name = ?)", *filters.DepartureCity)
	}
	if filters.ArrivalCity != nil {
		query = query.Where("ref_id IN (SELECT trip_ref_id FROM city_trips ct JOIN cities c ON c.ref_id = ct.city_ref_id WHERE ct.type = 'ARRIVAL' AND c.city_name = ?)", *filters.ArrivalCity)
	}
	if filters.Date != nil {
		query = query.Where("DATE(date_trip) = DATE(?)", *filters.Date)
	}

	var models []database.TripModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Trip, len(models))
	for i, m := range models {
		result[i] = toTripEntity(&m)
	}
	return result, nil
}

func (r *GormTripRepository) Create(ctx context.Context, data entities.CreateTripData) (*entities.Trip, error) {
	var trip *entities.Trip

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := &database.TripModel{
			DateTrip:    data.DateTrip,
			Kms:         data.Kms,
			Seats:       data.Seats,
			DriverRefID: data.DriverRefID,
			CarRefID:    data.CarRefID,
		}
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		// Create city-trip associations: first is DEPARTURE, second is ARRIVAL
		if len(data.CityRefIDs) >= 2 {
			cityTrips := []database.CityTripModel{
				{TripRefID: m.RefID, CityRefID: data.CityRefIDs[0], Type: "DEPARTURE"},
				{TripRefID: m.RefID, CityRefID: data.CityRefIDs[1], Type: "ARRIVAL"},
			}
			if err := tx.Create(&cityTrips).Error; err != nil {
				return err
			}
		}

		e := toTripEntity(m)
		trip = &e
		return nil
	})

	if err != nil {
		return nil, err
	}
	return trip, nil
}

func (r *GormTripRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m database.TripModel
		if err := tx.Where("id = ?", id).First(&m).Error; err != nil {
			return err
		}
		// Delete city-trip associations first
		if err := tx.Where("trip_ref_id = ?", m.RefID).Delete(&database.CityTripModel{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&database.TripModel{}).Error
	})
}

func toTripEntity(m *database.TripModel) entities.Trip {
	return entities.Trip{
		ID: m.ID, RefID: m.RefID,
		DateTrip: m.DateTrip, Kms: m.Kms, Seats: m.Seats,
		DriverRefID: m.DriverRefID, CarRefID: m.CarRefID,
	}
}
