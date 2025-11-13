-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE city_status AS ENUM (
    'supported',
    'suspended',
    'unsupported'
);

CREATE TABLE cities (
    id         UUID                  PRIMARY KEY NOT NULL,
    country_id VARCHAR(3)            NOT NULL,
    status     city_status           NOT NULL,
    point      geography(Point,4326) NOT NULL, -- lon/lat
    name       VARCHAR(255)          NOT NULL, -- default name in English
    icon       VARCHAR(255),
    slug       VARCHAR(255)          UNIQUE,
    timezone   VARCHAR(64)           NOT NULL, -- IANA tz

    created_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
);


-- +migrate Down
DROP TABLE IF EXISTS cities CASCADE;

DROP TYPE IF EXISTS city_status;