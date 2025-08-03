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
    "id"          UUID         PRIMARY KEY NOT NULL,
    "country_id"  UUID         NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "name"        VARCHAR(255) NOT NULL UNIQUE,
    "status"      VARCHAR(32)  NOT NULL,
    "created_at"  TIMESTAMP    NOT NULL,
    "updated_at"  TIMESTAMP    NOT NULL
);

CREATE TYPE city_admin_roles AS ENUM (
    'owner',
    'admin',
    'moderator'
);

CREATE TABLE "cities_admins" (
    "user_id"    UUID             NOT NULL,
    "city_id"    UUID             NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "role"       city_admin_roles NOT NULL,
    "created_at" TIMESTAMP        NOT NULL,

    PRIMARY KEY ("user_id", "city_id"),
);

CREATE UNIQUE INDEX city_owner_unique
    ON city_admins(city_id)
    WHERE role = 'owner';

-- +migrate Down
DROP TABLE IF EXISTS "cities_admins" CASCADE;
DROP TABLE IF EXISTS "city" CASCADE;
DROP TABLE IF EXISTS "countries" CASCADE;

DROP INDEX IF EXISTS city_owner_unique;

DROP TYPE IF EXISTS city_admin_roles;