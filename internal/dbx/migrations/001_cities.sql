-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE country_statuses AS ENUM (
    'supported',
    'unsupported',
);

CREATE TABLE "countries" (
    "id"         UUID              PRIMARY KEY NOT NULL,
    "name"       VARCHAR(255)      NOT NULL UNIQUE,
    "status"     country_statuses  NOT NULL,
    "created_at" TIMESTAMP         NOT NULL,
    "updated_at" TIMESTAMP         NOT NULL
);

CREATE TYPE city_statuses AS ENUM (
    'supported',
    'unsupported',
);

CREATE TABLE "city" (
    "id"         UUID                         PRIMARY KEY NOT NULL,
    "country_id" UUID                         NOT NULL REFERENCES "countries" ("id") ON DELETE CASCADE,
    "status"     city_statuses                NOT NULL,
    "center"     geometry(Point, 4326)        NOT NULL,
    "boundary"   geometry(MultiPolygon, 4326) NOT NULL,
    "icon"       VARCHAR(255)                 NOT NULL,
    "slug"       VARCHAR(255)                 NOT NULL UNIQUE,
    "timezone"   VARCHAR(64)                  NOT NULL, -- IANA tz (example: Europe/Kyiv)

    "created_at" TIMESTAMP                    NOT NULL,
    "updated_at" TIMESTAMP                    NOT NULL
);

CREATE TYPE "city_languages" AS ENUM (
    'en', 'es', 'fr', 'de', 'it', 'pt', 'uk'
);

CREATE TABLE "city_details" (
    "city_id"     UUID           NOT NULL REFERENCES "city" ("id") ON DELETE CASCADE,
    "language"    city_languages NOT NULL,
    "name"        VARCHAR(255)   NOT NULL,
    "description" VARCHAR(2047),

    PRIMARY KEY (city_id, language)
);

CREATE OR REPLACE FUNCTION city_guard_country_status()
RETURNS trigger AS $$
DECLARE
    cstatus country_statuses;
BEGIN
    SELECT status INTO cstatus FROM countries WHERE id = NEW.country_id;

    IF cstatus = 'unsupported' THEN
        IF TG_OP = 'INSERT' THEN
            RAISE EXCEPTION 'Cannot create city in unsupported country (%).', NEW.country_id;
        ELSIF TG_OP = 'UPDATE' THEN
            -- Разрешаем апдейты, кроме попытки сделать город supported
            IF NEW.status = 'supported' THEN
                RAISE EXCEPTION 'Cannot set city to supported while country (%) is unsupported.', NEW.country_id;
            END IF;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_city_insert_update_guard
BEFORE INSERT OR UPDATE ON city
FOR EACH ROW
EXECUTE FUNCTION city_guard_country_status();

-- 2) Если страна стала unsupported — все её города делаем unsupported.
CREATE OR REPLACE FUNCTION enforce_country_status_on_cities()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'unsupported' THEN
        UPDATE city
            SET status = 'unsupported'
        WHERE country_id = NEW.id
            AND status <> 'unsupported';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_countries_status_propagate
AFTER UPDATE OF status ON countries
FOR EACH ROW
WHEN (OLD.status IS DISTINCT FROM NEW.status)
EXECUTE FUNCTION enforce_country_status_on_cities();

CREATE OR REPLACE FUNCTION city_details_require_at_least_one()
RETURNS trigger AS $$
BEGIN
-- проверяем только если город ещё существует (не удаляется каскадно)
    IF EXISTS (SELECT 1 FROM city WHERE id = OLD.city_id) THEN
        IF NOT EXISTS (SELECT 1 FROM city_details WHERE city_id = OLD.city_id) THEN
            RAISE EXCEPTION 'Cannot delete the last city_details row for city %', OLD.city_id;
        END IF;
        END IF;
    RETURN NULL; -- AFTER/CONSTRAINT триггер значением не пользуется
END;
$$ LANGUAGE plpgsql;

-- деферрируемый constraint trigger:
-- - AFTER DELETE: запретить удаление последней детали
-- - AFTER UPDATE OF city_id: если переносим запись на другой city, у старого города не должно остаться 0
CREATE CONSTRAINT TRIGGER trg_city_details_nonempty
AFTER DELETE OR UPDATE OF city_id ON city_details
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION city_details_require_at_least_one();

-- +migrate Down
DROP TRIGGER IF EXISTS trg_countries_status_propagate ON countries;
DROP FUNCTION IF EXISTS enforce_country_status_on_cities();

DROP TRIGGER IF EXISTS trg_city_insert_update_guard ON city;
DROP FUNCTION IF EXISTS city_guard_country_status();

DROP TRIGGER IF EXISTS trg_city_details_nonempty ON city_details;
DROP FUNCTION IF EXISTS city_details_require_at_least_one();

DROP TABLE IF EXISTS "countries" CASCADE;
DROP TABLE IF EXISTS "city_details" CASCADE;
DROP TABLE IF EXISTS "city" CASCADE;

DROP TYPE IF EXISTS city_statuses;