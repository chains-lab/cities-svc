-- name: CreateCityGov :one
-- $user_id, $city_id, $role
INSERT INTO city_governments (user_id, city_id, role, created_at, updated_at)
VALUES (sqlc.arg(user_id), sqlc.arg(city_id), sqlc.arg(role), now(), now())
    RETURNING *;

-- name: DeleteCityAdmin :exec
-- $city_id
DELETE FROM city_governments
WHERE city_id = sqlc.arg(city_id)
  AND role = 'admin';

-- name: GetCityAdmin :one
-- $city_id
SELECT *
FROM city_governments
WHERE city_id = sqlc.arg(city_id)
  AND role = 'admin';

-- name: GetCityGov :one
-- $user_id, $city_id
SELECT *
FROM city_governments
WHERE user_id = sqlc.arg(user_id)
  AND city_id = sqlc.arg(city_id);

-- name: DeleteCityGov :exec
-- $user_id, $city_id
DELETE FROM city_governments
WHERE user_id = sqlc.arg(user_id)
  AND city_id = sqlc.arg(city_id);

-- name: SelectCityGovs :many
-- $city_id, $page, $page_size
-- Пагинация + сортировка: сначала admin, потом moderator; внутри — по created_at DESC
WITH base AS (
    SELECT
        cg.*,
        COUNT(*) OVER() AS total_count
    FROM city_governments cg
    WHERE cg.city_id = sqlc.arg(city_id)
)
SELECT *
FROM base
ORDER BY
    CASE WHEN role = 'admin' THEN 0 ELSE 1 END,
    created_at DESC
    LIMIT  sqlc.arg(page_size)::int8
OFFSET (GREATEST(sqlc.arg(page)::int8, 1) - 1) * sqlc.arg(page_size)::int8;
