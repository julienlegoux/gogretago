package brand

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListBrands_Success(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	brands := []entities.Brand{
		{ID: "brand-1", RefID: 1, Name: "Toyota"},
		{ID: "brand-2", RefID: 2, Name: "Honda"},
	}

	brandRepo := mocks.NewMockBrandRepository(t)

	brandRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(brands, 2, nil)

	uc := NewListBrandsUseCase(brandRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
}
