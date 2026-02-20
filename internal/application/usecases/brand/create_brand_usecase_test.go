package brand

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBrand_Success(t *testing.T) {
	ctx := context.Background()
	input := dtos.CreateBrandInput{Name: "Toyota"}

	expectedBrand := &entities.Brand{ID: "brand-1", RefID: 1, Name: "Toyota"}

	brandRepo := mocks.NewMockBrandRepository(t)

	brandRepo.EXPECT().Create(mock.Anything, entities.CreateBrandData{
		Name: "Toyota",
	}).Return(expectedBrand, nil)

	uc := NewCreateBrandUseCase(brandRepo)
	result, err := uc.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "brand-1", result.ID)
	assert.Equal(t, "Toyota", result.Name)
}

func TestCreateBrand_RepoError(t *testing.T) {
	ctx := context.Background()
	input := dtos.CreateBrandInput{Name: "Toyota"}

	brandRepo := mocks.NewMockBrandRepository(t)

	brandRepo.EXPECT().Create(mock.Anything, entities.CreateBrandData{
		Name: "Toyota",
	}).Return(nil, errors.New("database error"))

	uc := NewCreateBrandUseCase(brandRepo)
	result, err := uc.Execute(ctx, input)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}
