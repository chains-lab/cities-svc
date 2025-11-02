-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE city_statuses (
    id          VARCHAR(32) PRIMARY KEY CHECK ('^[a-z_]{1,32}$'),
    name        VARCHAR(64) NOT NULL,
    description VARCHAR(255) NOT NULL,
    accessible  BOOLEAN     NOT NULL DEFAULT TRUE,

    allowed_administration BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
)

CREATE TABLE city (
    id         UUID          PRIMARY KEY NOT NULL,
    country_id VARCHAR(3)    NOT NULL,
    status     VARCHAR(32)   NOT NULL REFERENCES city_statuses(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    name       VARCHAR(255)          NOT NULL,
    icon       VARCHAR(255),
    slug       VARCHAR(255)          UNIQUE,
    timezone   VARCHAR(64)           NOT NULL, -- IANA tz
    point      geography(Point,4326) NOT NULL, -- lon/lat

    created_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

-- +migrate Down
DROP TABLE IF EXISTS city CASCADE;
DROP TABLE IF EXISTS city_statuses CASCADE;