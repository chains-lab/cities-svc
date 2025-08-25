-- name: CreateCity :one
-- $1=id, $2=country_id, $3=status, $4=zone_geojson, $5=icon, $6=slug, $7=timezone
INSERT INTO city (id, country_id, status, zone, name, icon, slug, timezone, created_at, updated_at)
VALUES (
        uuid_generate_v4(),  -- $1
        $1,
        $2,
        ST_SetSRID(ST_GeomFromGeoJSON($3), 4326),  -- MultiPolygon GeoJSON
        $4,
        $5,
        $6,
        now(),
        now(),
        now()
       )
    RETURNING *;

-- name: GetCityByID :one
SELECT *
FROM city
WHERE id = $1;

-- name: GetCityBySlug :one
SELECT id
FROM city
WHERE slug = $1;

-- name: GetCityWithDetails :one
-- $1 = city_id (uuid), $2 = language (city_languages)
SELECT
    c.id,
    c.country_id,
    c.status,
    c.zone,
    c.icon,
    c.slug,
    c.timezone,
    c.created_at,
    c.updated_at,
    d.language,
    d.name
FROM city c
         JOIN city_details d
              ON d.city_id = c.id
WHERE c.id = $1
  AND d.language = $2;


-- name: GetNearestCity :one
WITH p AS (
    SELECT ST_SetSRID(ST_MakePoint($1, $2), 4326)::geometry(Point,4326) AS pt
)
SELECT c.id
FROM city c, p
ORDER BY
    CASE WHEN ST_Contains(c.zone, p.pt) THEN 0 ELSE 1 END,
    c.zone <-> p.pt    -- KNN, требует GIST-индекса на city.zone
LIMIT 1;

-- $1 = name_pattern (например '%roma%')
-- $2 = statuses (city_statuses[])  -- NULL или '{}' => без фильтра
-- $3 = country_ids (uuid[])        -- NULL или '{}' => без фильтра
-- $4 = page
-- $5 = size
-- name: SelectCityDetailsByNames :many
-- Возвращаем: city_id, language, name, total_count
WITH matched AS (
    SELECT cd.city_id, cd.language, cd.name
    FROM city_details cd
    WHERE cd.name ILIKE sqlc.arg(name_pattern)
    ),
    base AS (
SELECT
    cd.city_id,
    cd.language,
    cd.name,
    COUNT(*) OVER() AS total_count
FROM matched cd
    JOIN city c ON c.id = cd.city_id
WHERE
    (
    sqlc.narg(statuses)::city_statuses[] IS NULL
   OR cardinality(sqlc.narg(statuses)::city_statuses[]) = 0
   OR c.status = ANY(sqlc.narg(statuses)::city_statuses[])
    )
  AND
    (
    sqlc.narg(country_ids)::uuid[] IS NULL
   OR cardinality(sqlc.narg(country_ids)::uuid[]) = 0
   OR c.country_id = ANY(sqlc.narg(country_ids)::uuid[])
    )
    )
SELECT city_id, language, name, total_count
FROM base
ORDER BY name, city_id
LIMIT sqlc.arg(page_size)::int8
OFFSET (GREATEST(sqlc.arg(page)::int8, 1) - 1) * sqlc.arg(page_size)::int8;



-- name: UpdateCityStatus :exec
-- $1 = city_id (uuid), $2 = status (city_statuses)
UPDATE city
SET status = $2,
    updated_at = now()
WHERE id = $1;

-- name: UpdateCityCenter :exec
-- $1 = city_id (uuid), $2 = lon (float8), $3 = lat (float8)
UPDATE city
SET center = ST_SetSRID(ST_MakePoint($2, $3), 4326),
    updated_at = now()
WHERE id = $1;

-- name: UpdateCityZone :exec
-- $1 = city_id (uuid), $2 = zone (text, MultiPolygon in GeoJSON)
UPDATE city
SET zone   = ST_SetSRID(ST_GeomFromGeoJSON($2), 4326),
    updated_at = now()
WHERE id = $1;

-- name: UpdateCityIcon :exec
-- $1 = city_id (uuid), $2 = icon (varchar)
UPDATE city
SET icon = $2,
    updated_at = now()
WHERE id = $1;

-- name: UpdateCitySlug :exec
-- $1 = city_id (uuid), $2 = slug (varchar, UNIQUE)
UPDATE city
SET slug = $2,
    updated_at = now()
WHERE id = $1;

-- name: UpdateCityTimezone :exec
-- $1 = city_id (uuid), $2 = timezone (varchar, IANA)
UPDATE city
SET timezone = $2,
    updated_at = now()
WHERE id = $1;
