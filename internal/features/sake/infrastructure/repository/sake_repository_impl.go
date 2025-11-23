package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/database/sqlc"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// sakeRepositoryImpl 酒リポジトリの実装
type sakeRepositoryImpl struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewSakeRepository コンストラクタ
func NewSakeRepository(db *pgxpool.Pool) repository.SakeRepository {
	return &sakeRepositoryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}

// List 酒一覧を取得
func (r *sakeRepositoryImpl) List(ctx context.Context, filter repository.ListSakesFilter) ([]entity.Sake, entity.Pagination, error) {
	defer logger.TraceMethodAuto(ctx, filter)()

	// カウント取得
	countParams := sqlc.CountSakesParams{
		TypeID:    filter.TypeID,
		BreweryID: filter.BreweryID,
	}
	total, err := r.queries.CountSakes(ctx, countParams)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{
			"filter": filter,
		})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒の件数取得に失敗しました", err)
	}

	// リスト取得
	listParams := sqlc.ListSakesParams{
		Limit:     filter.Limit,
		Offset:    filter.Offset,
		TypeID:    filter.TypeID,
		BreweryID: filter.BreweryID,
	}
	sakeRows, err := r.queries.ListSakes(ctx, listParams)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{
			"filter": filter,
		})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧の取得に失敗しました", err)
	}

	// Entity変換
	sakes := make([]entity.Sake, 0, len(sakeRows))
	for _, row := range sakeRows {
		sake, err := r.toSakeEntity(ctx, row)
		if err != nil {
			return nil, entity.Pagination{}, err
		}
		sakes = append(sakes, *sake)
	}

	pagination := entity.Pagination{
		Total:  total,
		Offset: filter.Offset,
		Limit:  filter.Limit,
	}

	return sakes, pagination, nil
}

// toSakeEntity sqlcモデルからDomainエンティティに変換
func (r *sakeRepositoryImpl) toSakeEntity(ctx context.Context, row sqlc.Sake) (*entity.Sake, error) {
	// 酒の種類取得
	sakeType, err := r.queries.GetSakeType(ctx, row.TypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFoundError("酒の種類が見つかりません")
		}
		logger.LogDatabaseError(ctx, "SELECT", "sake_types", err, map[string]interface{}{
			"type_id": row.TypeID,
		})
		return nil, apperror.DatabaseError("酒の種類取得に失敗しました", err)
	}

	// 酒造取得
	brewery, err := r.queries.GetBrewery(ctx, row.BreweryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFoundError("酒造が見つかりません")
		}
		logger.LogDatabaseError(ctx, "SELECT", "breweries", err, map[string]interface{}{
			"brewery_id": row.BreweryID,
		})
		return nil, apperror.DatabaseError("酒造取得に失敗しました", err)
	}

	// 飲み方取得
	drinkStyleRows, err := r.queries.GetDrinkStylesBySakeID(ctx, row.ID)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "drink_styles", err, map[string]interface{}{
			"sake_id": row.ID,
		})
		return nil, apperror.DatabaseError("飲み方取得に失敗しました", err)
	}

	drinkStyles := make([]entity.DrinkStyle, 0, len(drinkStyleRows))
	for _, ds := range drinkStyleRows {
		drinkStyles = append(drinkStyles, entity.DrinkStyle{
			ID:          ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	// 座標抽出
	latitude, longitude := extractCoordinates(brewery.Position)

	return &entity.Sake{
		ID: row.ID,
		Type: entity.SakeType{
			ID:   sakeType.ID,
			Name: sakeType.Name,
		},
		Brewery: entity.Brewery{
			ID:            brewery.ID,
			Name:          brewery.Name,
			OriginCountry: brewery.OriginCountry,
			OriginRegion:  brewery.OriginRegion,
			Latitude:      latitude,
			Longitude:     longitude,
		},
		Name:        row.Name,
		ABV:         convertNumericToFloat32(row.Abv),
		TasteNotes:  row.TasteNotes,
		Memo:        row.Memo,
		DrinkStyles: drinkStyles,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

// extractCoordinates GEOMETRY型から緯度経度を抽出
func extractCoordinates(geomData interface{}) (*float64, *float64) {
	// PostGISのGEOMETRY型はinterface{}として返されるため
	// 現時点ではGEOMETRY型のパースは未実装(座標はnilを返す)
	// 将来的にはtwpayne/go-geomを使用して座標を抽出する
	return nil, nil
}

// convertNumericToFloat32 pgtype.NumericをFloat32に変換
func convertNumericToFloat32(n pgtype.Numeric) float32 {
	f64, err := n.Float64Value()
	if err != nil || !f64.Valid {
		return 0.0
	}
	return float32(f64.Float64)
}
