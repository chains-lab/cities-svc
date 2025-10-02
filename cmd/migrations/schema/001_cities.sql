-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE country_statuses AS ENUM (
    'supported',
    'deprecated',
    'unsupported'
);

CREATE TABLE countries (
    id         UUID             PRIMARY KEY NOT NULL,
    name       VARCHAR(255)     NOT NULL UNIQUE,
    status     country_statuses NOT NULL,

    created_at TIMESTAMP        NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP        NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TYPE city_statuses AS ENUM (
    'official',
    'community',
    'deprecated'
);

CREATE TABLE city (
    id         UUID                  PRIMARY KEY NOT NULL,
    country_id UUID                  NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    point      geography(Point,4326) NOT NULL, -- lon/lat
    status     city_statuses         NOT NULL,
    name       VARCHAR(255)          NOT NULL, -- default name in English
    icon       VARCHAR(255),
    slug       VARCHAR(255)          UNIQUE,
    timezone   VARCHAR(64)           NOT NULL, -- IANA tz

    created_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP             NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

-- +migrate Down
DROP TABLE IF EXISTS city CASCADE;
DROP TABLE IF EXISTS countries CASCADE;

DROP TYPE IF EXISTS city_statuses;
DROP TYPE IF EXISTS country_statuses;
