-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_gov_roles AS ENUM (
    'mayor',
    'government',
    'moderator'
);

CREATE TABLE "city_governments" (
    "user_id"    UUID           PRIMARY KEY, -- один юзер только в одном городе
    "city_id"    UUID           NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "role"       city_gov_roles NOT NULL,
    "label"      VARCHAR(255),
    "created_at" TIMESTAMP      NOT NULL,
    "updated_at" TIMESTAMP      NOT NULL
);

CREATE UNIQUE INDEX city_gov_unique
    ON city_governments(city_id)
    WHERE role = 'mayor';

-- +migrate Down
DROP INDEX IF EXISTS city_gov_unique;
DROP TABLE IF EXISTS "city_governments" CASCADE;
DROP TYPE IF EXISTS city_gov_roles;
DROP EXTENSION IF EXISTS "uuid-ossp";
