-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_gov_roles AS ENUM (
    'mayor',
    'advisors',
    'member',
    'moderator'
);

CREATE TABLE "city_governments" (
    "id"         UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id"    UUID           NOT NULL,
    "city_id"    UUID           NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "active"     BOOLEAN        NOT NULL DEFAULT TRUE,
    "role"       city_gov_roles NOT NULL,
    "label"      VARCHAR(255),
    "created_at" TIMESTAMP      NOT NULL,
    "updated_at" TIMESTAMP      NOT NULL,
    
    UNIQUE(user_id, city_id)
);

CREATE UNIQUE INDEX city_gov_unique
    ON city_governments(city_id)
    WHERE role = 'mayor' AND active = TRUE;

CREATE UNIQUE INDEX city_gov_user_unique
    ON city_governments(user_id)
    WHERE active = TRUE;

-- +migrate Down
DROP INDEX IF EXISTS city_gov_unique;
DROP INDEX IF EXISTS city_gov_user_unique;
DROP TABLE IF EXISTS "city_governments" CASCADE;
DROP TYPE IF EXISTS city_gov_roles;
DROP EXTENSION IF EXISTS "uuid-ossp";
