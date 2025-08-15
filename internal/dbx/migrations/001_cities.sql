-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "countries" (
    "id"         UUID         PRIMARY KEY NOT NULL,
    "name"       VARCHAR(255) NOT NULL UNIQUE,
    "status"     VARCHAR(32)  NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL
);

CREATE TYPE city_statuses AS ENUM (
    'supported',
    'suspended',
    'unsupported'
);

CREATE TABLE "city" (
    "id"          UUID          PRIMARY KEY NOT NULL,
    "country_id"  UUID          NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "name"        VARCHAR(255)  NOT NULL,
    "status"      city_statuses NOT NULL,
    "created_at"  TIMESTAMP     NOT NULL,
    "updated_at"  TIMESTAMP     NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "city" CASCADE;
DROP TABLE IF EXISTS "countries" CASCADE;

DROP TYPE IF EXISTS city_statuses;