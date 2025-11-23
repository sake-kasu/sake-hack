package entity

import "time"

// SakeType 酒の種類
type SakeType struct {
	ID   int32
	Name string
}

// Brewery 酒造
type Brewery struct {
	ID            int32
	Name          string
	OriginCountry string
	OriginRegion  *string
	Latitude      *float64
	Longitude     *float64
}

// DrinkStyle 飲み方
type DrinkStyle struct {
	ID          int32
	Name        string
	Description *string
}

// Sake 酒
type Sake struct {
	ID          int32
	Type        SakeType
	Brewery     Brewery
	Name        string
	ABV         float32
	TasteNotes  string
	Memo        *string
	DrinkStyles []DrinkStyle
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Pagination ページネーション情報
type Pagination struct {
	Total  int64
	Offset int32
	Limit  int32
}
