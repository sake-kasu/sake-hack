-- sake_typesテーブル
CREATE TABLE sake_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
COMMENT ON TABLE sake_types IS '酒の種類';
COMMENT ON COLUMN sake_types.name IS '酒の種類の名前';
COMMENT ON COLUMN sake_types.created_at IS '作成日時';
COMMENT ON COLUMN sake_types.updated_at IS '更新日時';

-- breweriesテーブル
CREATE TABLE breweries (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    origin_country TEXT NOT NULL,
    origin_region TEXT,
    position GEOMETRY(Point, 4326),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(name, origin_country)
);
COMMENT ON TABLE breweries IS '酒造';
COMMENT ON COLUMN breweries.name IS '酒造名';
COMMENT ON COLUMN breweries.origin_country IS '所在国';
COMMENT ON COLUMN breweries.origin_region IS '所在地域';
COMMENT ON COLUMN breweries.position IS '座標(任意)';
COMMENT ON COLUMN breweries.created_at IS '作成日時';
COMMENT ON COLUMN breweries.updated_at IS '更新日時';

-- インデックス作成
CREATE INDEX idx_breweries_origin_country ON breweries(origin_country);
CREATE INDEX idx_breweries_position ON breweries USING GIST(position);

-- sakesテーブル
CREATE TABLE sakes (
    id SERIAL PRIMARY KEY,
    type_id INTEGER NOT NULL REFERENCES sake_types(id) ON DELETE RESTRICT,
    brewery_id INTEGER NOT NULL REFERENCES breweries(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    abv NUMERIC(4, 2) NOT NULL,
    taste_notes TEXT NOT NULL,
    memo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
COMMENT ON TABLE sakes IS '酒';
COMMENT ON COLUMN sakes.type_id IS '酒の種類id';
COMMENT ON COLUMN sakes.brewery_id IS '酒造id';
COMMENT ON COLUMN sakes.name IS '酒名';
COMMENT ON COLUMN sakes.abv IS 'アルコール度数(%)';
COMMENT ON COLUMN sakes.taste_notes IS '味の特徴';
COMMENT ON COLUMN sakes.memo IS '感想';
COMMENT ON COLUMN sakes.created_at IS '作成日時';
COMMENT ON COLUMN sakes.updated_at IS '更新日時';

-- インデックス作成
CREATE INDEX idx_sakes_type_id ON sakes(type_id);
CREATE INDEX idx_sakes_brewery_id ON sakes(brewery_id);
CREATE INDEX idx_sakes_name ON sakes(name);

-- drink_stylesテーブル
CREATE TABLE drink_styles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
COMMENT ON TABLE drink_styles IS '飲み方(マスター)';
COMMENT ON COLUMN drink_styles.name IS '酒の飲み方';
COMMENT ON COLUMN drink_styles.description IS '詳細説明';
COMMENT ON COLUMN drink_styles.created_at IS '作成日時';
COMMENT ON COLUMN drink_styles.updated_at IS '更新日時';

-- sake_drink_stylesテーブル(中間テーブル)
CREATE TABLE sake_drink_styles (
    sake_id INTEGER NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    drink_style_id INTEGER NOT NULL REFERENCES drink_styles(id) ON DELETE RESTRICT,
    PRIMARY KEY (sake_id, drink_style_id)
);
COMMENT ON TABLE sake_drink_styles IS '酒-飲み方(中間テーブル)';
COMMENT ON COLUMN sake_drink_styles.sake_id IS '酒id';
COMMENT ON COLUMN sake_drink_styles.drink_style_id IS '酒の飲み方id';

-- インデックス作成
CREATE INDEX idx_sake_drink_styles_drink_style_id ON sake_drink_styles(drink_style_id);
