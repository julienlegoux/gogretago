package di

import (
	"github.com/lgxju/gogretago/internal/application/usecases/auth"
	"github.com/lgxju/gogretago/internal/application/usecases/brand"
	"github.com/lgxju/gogretago/internal/application/usecases/car"
	"github.com/lgxju/gogretago/internal/application/usecases/city"
	"github.com/lgxju/gogretago/internal/application/usecases/color"
	"github.com/lgxju/gogretago/internal/application/usecases/driver"
	"github.com/lgxju/gogretago/internal/application/usecases/inscription"
	"github.com/lgxju/gogretago/internal/application/usecases/trip"
	"github.com/lgxju/gogretago/internal/application/usecases/user"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/lgxju/gogretago/internal/infrastructure/database"
	infrarepos "github.com/lgxju/gogretago/internal/infrastructure/repositories"
	infraservices "github.com/lgxju/gogretago/internal/infrastructure/services"
	"gorm.io/gorm"
)

// Container holds all dependencies for the application
type Container struct {
	// Database
	DB *gorm.DB

	// Repositories
	AuthRepository        repositories.AuthRepository
	UserRepository        repositories.UserRepository
	DriverRepository      repositories.DriverRepository
	BrandRepository       repositories.BrandRepository
	ModelRepository       repositories.ModelRepository
	ColorRepository       repositories.ColorRepository
	CarRepository         repositories.CarRepository
	CityRepository        repositories.CityRepository
	TripRepository        repositories.TripRepository
	InscriptionRepository repositories.InscriptionRepository

	// Services
	PasswordService services.PasswordService
	JwtService      services.JwtService
	EmailService    services.EmailService

	// Auth Use Cases
	RegisterUseCase *auth.RegisterUseCase
	LoginUseCase    *auth.LoginUseCase

	// User Use Cases
	ListUsersUseCase     *user.ListUsersUseCase
	GetUserUseCase       *user.GetUserUseCase
	UpdateUserUseCase    *user.UpdateUserUseCase
	AnonymizeUserUseCase *user.AnonymizeUserUseCase

	// Driver Use Cases
	CreateDriverUseCase *driver.CreateDriverUseCase

	// Brand Use Cases
	ListBrandsUseCase  *brand.ListBrandsUseCase
	CreateBrandUseCase *brand.CreateBrandUseCase
	DeleteBrandUseCase *brand.DeleteBrandUseCase

	// Color Use Cases
	ListColorsUseCase  *color.ListColorsUseCase
	CreateColorUseCase *color.CreateColorUseCase
	UpdateColorUseCase *color.UpdateColorUseCase
	DeleteColorUseCase *color.DeleteColorUseCase

	// City Use Cases
	ListCitiesUseCase *city.ListCitiesUseCase
	CreateCityUseCase *city.CreateCityUseCase
	DeleteCityUseCase *city.DeleteCityUseCase

	// Car Use Cases
	ListCarsUseCase  *car.ListCarsUseCase
	CreateCarUseCase *car.CreateCarUseCase
	UpdateCarUseCase *car.UpdateCarUseCase
	DeleteCarUseCase *car.DeleteCarUseCase

	// Trip Use Cases
	ListTripsUseCase  *trip.ListTripsUseCase
	GetTripUseCase    *trip.GetTripUseCase
	FindTripsUseCase  *trip.FindTripsUseCase
	CreateTripUseCase *trip.CreateTripUseCase
	DeleteTripUseCase *trip.DeleteTripUseCase

	// Inscription Use Cases
	ListInscriptionsUseCase     *inscription.ListInscriptionsUseCase
	CreateInscriptionUseCase    *inscription.CreateInscriptionUseCase
	DeleteInscriptionUseCase    *inscription.DeleteInscriptionUseCase
	ListUserInscriptionsUseCase *inscription.ListUserInscriptionsUseCase
	ListTripPassengersUseCase   *inscription.ListTripPassengersUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer() (*Container, error) {
	// Database connection
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		return nil, err
	}

	// Create repositories
	authRepository := infrarepos.NewGormAuthRepository(db)
	userRepository := infrarepos.NewGormUserRepository(db)
	driverRepository := infrarepos.NewGormDriverRepository(db)
	brandRepository := infrarepos.NewGormBrandRepository(db)
	modelRepository := infrarepos.NewGormModelRepository(db)
	colorRepository := infrarepos.NewGormColorRepository(db)
	carRepository := infrarepos.NewGormCarRepository(db)
	cityRepository := infrarepos.NewGormCityRepository(db)
	tripRepository := infrarepos.NewGormTripRepository(db)
	inscriptionRepository := infrarepos.NewGormInscriptionRepository(db)

	// Create services
	passwordService := infraservices.NewArgonPasswordService()
	jwtService := infraservices.NewJwtService()
	emailService := infraservices.NewResendEmailService()

	// Auth use cases
	registerUseCase := auth.NewRegisterUseCase(authRepository, passwordService, emailService, jwtService)
	loginUseCase := auth.NewLoginUseCase(authRepository, userRepository, passwordService, jwtService)

	// User use cases
	listUsersUseCase := user.NewListUsersUseCase(userRepository)
	getUserUseCase := user.NewGetUserUseCase(userRepository)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepository)
	anonymizeUserUseCase := user.NewAnonymizeUserUseCase(userRepository)

	// Driver use cases
	createDriverUseCase := driver.NewCreateDriverUseCase(driverRepository, userRepository, authRepository)

	// Brand use cases
	listBrandsUseCase := brand.NewListBrandsUseCase(brandRepository)
	createBrandUseCase := brand.NewCreateBrandUseCase(brandRepository)
	deleteBrandUseCase := brand.NewDeleteBrandUseCase(brandRepository)

	// Color use cases
	listColorsUseCase := color.NewListColorsUseCase(colorRepository)
	createColorUseCase := color.NewCreateColorUseCase(colorRepository)
	updateColorUseCase := color.NewUpdateColorUseCase(colorRepository)
	deleteColorUseCase := color.NewDeleteColorUseCase(colorRepository)

	// City use cases
	listCitiesUseCase := city.NewListCitiesUseCase(cityRepository)
	createCityUseCase := city.NewCreateCityUseCase(cityRepository)
	deleteCityUseCase := city.NewDeleteCityUseCase(cityRepository)

	// Car use cases
	listCarsUseCase := car.NewListCarsUseCase(carRepository)
	createCarUseCase := car.NewCreateCarUseCase(carRepository, modelRepository, brandRepository, driverRepository)
	updateCarUseCase := car.NewUpdateCarUseCase(carRepository, modelRepository, brandRepository, driverRepository)
	deleteCarUseCase := car.NewDeleteCarUseCase(carRepository, driverRepository)

	// Trip use cases
	listTripsUseCase := trip.NewListTripsUseCase(tripRepository)
	getTripUseCase := trip.NewGetTripUseCase(tripRepository)
	findTripsUseCase := trip.NewFindTripsUseCase(tripRepository)
	createTripUseCase := trip.NewCreateTripUseCase(tripRepository, driverRepository, carRepository, cityRepository)
	deleteTripUseCase := trip.NewDeleteTripUseCase(tripRepository, driverRepository)

	// Inscription use cases
	listInscriptionsUseCase := inscription.NewListInscriptionsUseCase(inscriptionRepository)
	createInscriptionUseCase := inscription.NewCreateInscriptionUseCase(inscriptionRepository, userRepository, tripRepository)
	deleteInscriptionUseCase := inscription.NewDeleteInscriptionUseCase(inscriptionRepository)
	listUserInscriptionsUseCase := inscription.NewListUserInscriptionsUseCase(inscriptionRepository)
	listTripPassengersUseCase := inscription.NewListTripPassengersUseCase(inscriptionRepository)

	return &Container{
		DB: db,

		// Repositories
		AuthRepository:        authRepository,
		UserRepository:        userRepository,
		DriverRepository:      driverRepository,
		BrandRepository:       brandRepository,
		ModelRepository:       modelRepository,
		ColorRepository:       colorRepository,
		CarRepository:         carRepository,
		CityRepository:        cityRepository,
		TripRepository:        tripRepository,
		InscriptionRepository: inscriptionRepository,

		// Services
		PasswordService: passwordService,
		JwtService:      jwtService,
		EmailService:    emailService,

		// Auth
		RegisterUseCase: registerUseCase,
		LoginUseCase:    loginUseCase,

		// User
		ListUsersUseCase:     listUsersUseCase,
		GetUserUseCase:       getUserUseCase,
		UpdateUserUseCase:    updateUserUseCase,
		AnonymizeUserUseCase: anonymizeUserUseCase,

		// Driver
		CreateDriverUseCase: createDriverUseCase,

		// Brand
		ListBrandsUseCase:  listBrandsUseCase,
		CreateBrandUseCase: createBrandUseCase,
		DeleteBrandUseCase: deleteBrandUseCase,

		// Color
		ListColorsUseCase:  listColorsUseCase,
		CreateColorUseCase: createColorUseCase,
		UpdateColorUseCase: updateColorUseCase,
		DeleteColorUseCase: deleteColorUseCase,

		// City
		ListCitiesUseCase: listCitiesUseCase,
		CreateCityUseCase: createCityUseCase,
		DeleteCityUseCase: deleteCityUseCase,

		// Car
		ListCarsUseCase:  listCarsUseCase,
		CreateCarUseCase: createCarUseCase,
		UpdateCarUseCase: updateCarUseCase,
		DeleteCarUseCase: deleteCarUseCase,

		// Trip
		ListTripsUseCase:  listTripsUseCase,
		GetTripUseCase:    getTripUseCase,
		FindTripsUseCase:  findTripsUseCase,
		CreateTripUseCase: createTripUseCase,
		DeleteTripUseCase: deleteTripUseCase,

		// Inscription
		ListInscriptionsUseCase:     listInscriptionsUseCase,
		CreateInscriptionUseCase:    createInscriptionUseCase,
		DeleteInscriptionUseCase:    deleteInscriptionUseCase,
		ListUserInscriptionsUseCase: listUserInscriptionsUseCase,
		ListTripPassengersUseCase:   listTripPassengersUseCase,
	}, nil
}
