-- name: ListSakes :many
SELECT
    s.id,
    s.type_id,
    s.brewery_id,
    s.name,
    s.abv,
    s.taste_notes,
    s.memo,
    s.created_at,
    s.updated_at
FROM sakes s
WHERE
    (sqlc.narg('type_id')::INTEGER IS NULL OR s.type_id = sqlc.narg('type_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'))
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSakes :one
SELECT COUNT(*) AS total
FROM sakes s
WHERE
    (sqlc.narg('type_id')::INTEGER IS NULL OR s.type_id = sqlc.narg('type_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'));

-- name: GetSakeType :one
SELECT id, name, created_at, updated_at
FROM sake_types
WHERE id = $1;

-- name: GetBrewery :one
SELECT id, name, origin_country, origin_region, position, created_at, updated_at
FROM breweries
WHERE id = $1;

-- name: GetDrinkStylesBySakeID :many
SELECT
    ds.id,
    ds.name,
    ds.description,
    ds.created_at,
    ds.updated_at
FROM drink_styles ds
INNER JOIN sake_drink_styles sds ON ds.id = sds.drink_style_id
WHERE sds.sake_id = $1
ORDER BY ds.id;
