-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "countries" (
    "id"         UUID         PRIMARY KEY NOT NULL,
    "name"       VARCHAR(255) NOT NULL,
    "status"     VARCHAR(32)  NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL
);

CREATE TABLE "city" (
    "id"         UUID         PRIMARY KEY NOT NULL,
    "country_id" UUID         NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "name"       VARCHAR(255) NOT NULL,
    "status"     VARCHAR(32)  NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL
);

CREATE TABLE "city_admins" (
    "user_id"    UUID      NOT NULL,
    "city_id"    UUID      NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "city_admins" CASCADE;
DROP TABLE IF EXISTS "city" CASCADE;
DROP TABLE IF EXISTS "countries" CASCADE;