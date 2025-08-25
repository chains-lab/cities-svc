-- name: CreateCountry :one
-- $1 = name, $2 = status
INSERT INTO countries (id, name, status, created_at, updated_at)
VALUES ( uuid_generate_v4(), $1, $2, now(), now())
    RETURNING *;

-- name: GetCountryByID :one
SELECT *
FROM countries
WHERE id = $1;

-- name: GetCountryByName :one
SELECT *
FROM countries
WHERE name = $1;

-- name: UpdateCountryStatus :one
-- $1 = country_id, $2 = status
UPDATE countries
SET status = $2,
    updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: UpdateCountryName :one
-- $1 = country_id, $2 = new name
UPDATE countries
SET name = $2,
    updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: SearchBeNameAndStatuses :many
-- Ищем страны по имени + (опционально) по массиву статусов
WITH base AS (
    SELECT
        c.*,
        COUNT(*) OVER() AS total_count
    FROM countries c
    WHERE
        c.name ILIKE sqlc.arg(name_pattern)
    AND (
    sqlc.narg(statuses)::country_statuses[] IS NULL
    OR cardinality(sqlc.narg(statuses)::country_statuses[]) = 0
    OR c.status = ANY(sqlc.narg(statuses)::country_statuses[])
    )
    )
SELECT *
FROM base
ORDER BY created_at DESC, id
    LIMIT  sqlc.arg(page_size)::int8
OFFSET (GREATEST(sqlc.arg(page)::int8, 1) - 1) * sqlc.arg(page_size)::int8;
