package repository

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// ListSakesFilter 酒一覧取得のフィルター条件
type ListSakesFilter struct {
	TypeID    *int32
	BreweryID *int32
	Offset    int32
	Limit     int32
}

// SakeRepository 酒リポジトリのインターフェース
type SakeRepository interface {
	// List フィルター条件に基づいて酒一覧を取得
	List(ctx context.Context, filter ListSakesFilter) ([]entity.Sake, entity.Pagination, error)
}
