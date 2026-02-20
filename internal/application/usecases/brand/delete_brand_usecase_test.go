package brand

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteBrand_Success(t *testing.T) {
	ctx := context.Background()
	brandID := "brand-1"

	existing := &entities.Brand{ID: brandID, RefID: 1, Name: "Toyota"}

	brandRepo := mocks.NewMockBrandRepository(t)

	brandRepo.EXPECT().FindByID(mock.Anything, brandID).Return(existing, nil)
	brandRepo.EXPECT().Delete(mock.Anything, brandID).Return(nil)

	uc := NewDeleteBrandUseCase(brandRepo)
	err := uc.Execute(ctx, brandID)

	assert.NoError(t, err)
}

func TestDeleteBrand_NotFound(t *testing.T) {
	ctx := context.Background()
	brandID := "brand-nonexistent"

	brandRepo := mocks.NewMockBrandRepository(t)

	brandRepo.EXPECT().FindByID(mock.Anything, brandID).Return(nil, nil)

	uc := NewDeleteBrandUseCase(brandRepo)
	err := uc.Execute(ctx, brandID)

	assert.Error(t, err)
	var notFoundErr *domainerrors.BrandNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}
