package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
)

// MockSakeRepository はSakeRepositoryのモック
type MockSakeRepository struct {
	mock.Mock
}

func (m *MockSakeRepository) List(ctx context.Context, filter repository.ListSakesFilter) ([]entity.Sake, entity.Pagination, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, entity.Pagination{}, args.Error(2)
	}
	return args.Get(0).([]entity.Sake), args.Get(1).(entity.Pagination), args.Error(2)
}

func TestListSakesUsecase_Execute_Success(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	typeID := int32(1)
	breweryID := int32(2)
	expectedSakes := []entity.Sake{
		{
			ID:   1,
			Name: "獺祭",
		},
	}
	expectedPagination := entity.Pagination{
		Total:  100,
		Offset: 0,
		Limit:  20,
	}

	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ListSakesFilter) bool {
		return filter.Offset == 0 && filter.Limit == 20 &&
			filter.TypeID != nil && *filter.TypeID == typeID &&
			filter.BreweryID != nil && *filter.BreweryID == breweryID
	})).Return(expectedSakes, expectedPagination, nil)

	input := ListSakesInput{
		TypeID:    &typeID,
		BreweryID: &breweryID,
		Offset:    0,
		Limit:     20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Sakes, 1)
	assert.Equal(t, int32(1), output.Sakes[0].ID)
	assert.Equal(t, int64(100), output.Pagination.Total)
	mockRepo.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_OffsetValidation(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ListSakesFilter) bool {
		return filter.Offset == 0 // デフォルト値
	})).Return([]entity.Sake{}, entity.Pagination{}, nil)

	// offset < 0 の場合、デフォルト値0が使用される
	input := ListSakesInput{
		Offset: -1,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_LimitValidation_LessThan1(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ListSakesFilter) bool {
		return filter.Limit == 20 // デフォルト値
	})).Return([]entity.Sake{}, entity.Pagination{}, nil)

	// limit < 1 の場合、デフォルト値20が使用される
	input := ListSakesInput{
		Offset: 0,
		Limit:  0,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_LimitValidation_GreaterThan100(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ListSakesFilter) bool {
		return filter.Limit == 20 // デフォルト値
	})).Return([]entity.Sake{}, entity.Pagination{}, nil)

	// limit > 100 の場合、デフォルト値20が使用される
	input := ListSakesInput{
		Offset: 0,
		Limit:  101,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_RepositoryError(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	expectedErr := errors.New("database error")
	mockRepo.On("List", mock.Anything, mock.Anything).Return(nil, entity.Pagination{}, expectedErr)

	input := ListSakesInput{
		Offset: 0,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_EmptyResult(t *testing.T) {
	mockRepo := new(MockSakeRepository)
	uc := NewListSakesUsecase(mockRepo)

	mockRepo.On("List", mock.Anything, mock.Anything).Return([]entity.Sake{}, entity.Pagination{
		Total:  0,
		Offset: 0,
		Limit:  20,
	}, nil)

	input := ListSakesInput{
		Offset: 0,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Sakes, 0)
	assert.Equal(t, int64(0), output.Pagination.Total)
	mockRepo.AssertExpectations(t)
}
