package di

import (
	"github.com/lgxju/gogretago/internal/application/usecases/auth"
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
	UserRepository repositories.UserRepository

	// Services
	PasswordService services.PasswordService
	JwtService      services.JwtService
	EmailService    services.EmailService

	// Use Cases
	RegisterUseCase *auth.RegisterUseCase
	LoginUseCase    *auth.LoginUseCase
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
	userRepository := infrarepos.NewGormUserRepository(db)

	// Create services
	passwordService := infraservices.NewArgonPasswordService()
	jwtService := infraservices.NewJwtService()
	emailService := infraservices.NewResendEmailService()

	// Create use cases with injected dependencies
	registerUseCase := auth.NewRegisterUseCase(
		userRepository,
		passwordService,
		emailService,
		jwtService,
	)

	loginUseCase := auth.NewLoginUseCase(
		userRepository,
		passwordService,
		jwtService,
	)

	return &Container{
		DB:              db,
		UserRepository:  userRepository,
		PasswordService: passwordService,
		JwtService:      jwtService,
		EmailService:    emailService,
		RegisterUseCase: registerUseCase,
		LoginUseCase:    loginUseCase,
	}, nil
}
