-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE country_statuses AS ENUM (
    'supported',
    'unsupported'
);

CREATE TABLE "countries" (
    "id"         UUID              PRIMARY KEY NOT NULL ,
    "name"       VARCHAR(255)      NOT NULL UNIQUE,
    "status"     country_statuses  NOT NULL,
    "created_at" TIMESTAMP         NOT NULL,
    "updated_at" TIMESTAMP         NOT NULL
);

CREATE TYPE city_statuses AS ENUM (
    'supported',
    'unsupported'
);

CREATE TABLE "city" (
    "id"         UUID                         PRIMARY KEY NOT NULL,
    "country_id" UUID                         NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "status"     city_statuses                NOT NULL,
    "zone"       geometry(MultiPolygon, 4326) NOT NULL,
    "name"       VARCHAR(255)                 NOT NULL, -- default name in English
    "icon"       VARCHAR(255)                 NOT NULL,
    "slug"       VARCHAR(255)                 NOT NULL UNIQUE,
    "timezone"   VARCHAR(64)                  NOT NULL, -- IANA tz (example: Europe/Kyiv)

    "created_at" TIMESTAMP                    NOT NULL,
    "updated_at" TIMESTAMP                    NOT NULL
);

CREATE TYPE "city_languages" AS ENUM (
    'es', 'fr', 'de', 'it', 'pt', 'uk'
);

CREATE TABLE "city_details" (
    "city_id"     UUID           NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "language"    city_languages NOT NULL,
    "name"        VARCHAR(255)   NOT NULL,

    PRIMARY KEY (city_id, language)
);

-- +migrate Down
DROP TABLE IF EXISTS "countries" CASCADE;
DROP TABLE IF EXISTS "city_details" CASCADE;
DROP TABLE IF EXISTS "city" CASCADE;

DROP TYPE IF EXISTS city_statuses;