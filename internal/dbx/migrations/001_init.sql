-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "countries" (
    "id"         UUID         PRIMARY KEY NOT NULL,
    "name"       VARCHAR(255) NOT NULL UNIQUE,
    "status"     VARCHAR(32)  NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL
);

CREATE TABLE "city" (
    "id"          UUID         PRIMARY KEY NOT NULL,
    "country_id"  UUID         NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "name"        VARCHAR(255) NOT NULL,
    "status"      VARCHAR(32)  NOT NULL,
    "created_at"  TIMESTAMP    NOT NULL,
    "updated_at"  TIMESTAMP    NOT NULL
);

CREATE TYPE city_admin_roles AS ENUM (
    'admin',
    'moderator'
);

CREATE TABLE "cities_governments" (
    "user_id"    UUID             NOT NULL UNIQUE ,
    "city_id"    UUID             NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "role"       city_admin_roles NOT NULL,
    "created_at" TIMESTAMP        NOT NULL,
    "updated_at" TIMESTAMP        NOT NULL,

    PRIMARY KEY ("user_id", "city_id"),
);

CREATE UNIQUE INDEX city_owner_unique
    ON city_admins(city_id)
    WHERE role = 'admin';

CREATE TYPE forms_to_create_city_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

CREATE TABLE "forms_to_create_city" (
    "id"               UUID          PRIMARY KEY NOT NULL,
    "status"           VARCHAR(32)   forms_to_create_city_status NOT NULL DEFAULT 'pending',
    "initiator_id"     UUID          NOT NULL,
    "city_name"        VARCHAR(255)  NOT NULL,
    "country_id"       UUID          NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "contact_email"    VARCHAR(255)  NOT NULL,
    "contact_phone"    VARCHAR(50)   NOT NULL,
    "text"             VARCHAR(5096) NOT NULL,
    "user_reviewed_id" UUID          DEFAULT "00000000-0000-0000-0000-000000000000",
    "create_at"        TIMESTAMP     NOT NULL,
);

-- +migrate Down
DROP TABLE IF EXISTS "cities_admins" CASCADE;
DROP TABLE IF EXISTS "city" CASCADE;
DROP TABLE IF EXISTS "countries" CASCADE;

DROP INDEX IF EXISTS city_owner_unique;

DROP TYPE IF EXISTS city_admin_roles;