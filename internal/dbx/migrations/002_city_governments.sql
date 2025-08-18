-- +migrate Up

CREATE TYPE city_admin_roles AS ENUM (
    'admin',
    'moderator'
);

CREATE TABLE "city_governments" (
    "user_id"    UUID             NOT NULL UNIQUE ,
    "city_id"    UUID             NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "role"       city_admin_roles NOT NULL,
    "created_at" TIMESTAMP        NOT NULL,
    "updated_at" TIMESTAMP        NOT NULL,

    PRIMARY KEY ("user_id", "city_id"),
);

CREATE UNIQUE INDEX city_admin_unique
    ON city_admins(city_id)
    WHERE role = 'admin';

-- +migrate Down
DROP TABLE IF EXISTS "city_governments" CASCADE;
DROP INDEX IF EXISTS city_admin_unique;

DROP TYPE IF EXISTS city_admin_roles;