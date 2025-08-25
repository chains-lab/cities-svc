-- name: CreateCityDetails :exec
-- $1 = city_id (uuid), $2 = language (city_languages), $3 = name (varchar)
INSERT INTO city_details (city_id, language, name)
VALUES ($1, $2, $3);

-- name: GetCityDetailsByCityIdAndLanguage :one
-- $1 = city_id (uuid), $2 = language (city_languages)
SELECT city_id, language, name
FROM city_details
WHERE city_id = $1
  AND language = $2;

-- name: GetCityDetailsByCityIdAndAnyLanguage :one
-- $1 = city_id (uuid)
SELECT city_id, language, name
FROM city_details
WHERE city_id = $1
LIMIT 1;

-- name: SelectCityDetailsByCityId :many
-- $1 = city_id (uuid), $2 = page, $3 = size
WITH base AS (
    SELECT
        cd.city_id,
        cd.language,
        cd.name,
        COUNT(*) OVER() AS total_count
    FROM city_details cd
    WHERE cd.city_id = $1
)
SELECT *
FROM base
ORDER BY language
    LIMIT $3
OFFSET (GREATEST($2, 1) - 1) * $3;

-- name: UpdateCityDetails :exec
-- $1 = city_id (uuid), $2 = language (city_languages), $3 = new name (varchar)
UPDATE city_details
SET name = $3
WHERE city_id = $1
  AND language = $2;

-- name: DeleteCityDetails :exec
-- $1 = city_id (uuid), $2 = language (city_languages)
DELETE FROM city_details
WHERE city_id = $1
  AND language = $2;

