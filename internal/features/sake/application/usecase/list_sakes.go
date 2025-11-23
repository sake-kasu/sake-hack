package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ListSakesInput 酒一覧取得の入力パラメータ
type ListSakesInput struct {
	TypeID    *int32
	BreweryID *int32
	Offset    int32
	Limit     int32
}

// ListSakesOutput 酒一覧取得の出力
type ListSakesOutput struct {
	Sakes      []entity.Sake
	Pagination entity.Pagination
}

// ListSakesUsecaseInterface 酒一覧取得ユースケースのインターフェイス
type ListSakesUsecaseInterface interface {
	Execute(ctx context.Context, input ListSakesInput) (*ListSakesOutput, error)
}

// ListSakesUsecase 酒一覧取得ユースケース
type ListSakesUsecase struct {
	sakeRepo repository.SakeRepository
}

// NewListSakesUsecase コンストラクタ
func NewListSakesUsecase(sakeRepo repository.SakeRepository) *ListSakesUsecase {
	return &ListSakesUsecase{
		sakeRepo: sakeRepo,
	}
}

// Execute 酒一覧を取得する
func (u *ListSakesUsecase) Execute(ctx context.Context, input ListSakesInput) (*ListSakesOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	// バリデーション
	if input.Offset < 0 {
		input.Offset = 0
	}
	if input.Limit < 1 || input.Limit > 100 {
		input.Limit = 20
	}

	// リポジトリから取得
	sakes, pagination, err := u.sakeRepo.List(ctx, repository.ListSakesFilter{
		TypeID:    input.TypeID,
		BreweryID: input.BreweryID,
		Offset:    input.Offset,
		Limit:     input.Limit,
	})
	if err != nil {
		return nil, err
	}

	return &ListSakesOutput{
		Sakes:      sakes,
		Pagination: pagination,
	}, nil
}
