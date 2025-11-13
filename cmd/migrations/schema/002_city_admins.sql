-- +migrate Up
CREATE TYPE city_admin_role AS ENUM (
    'chief',
    'vice-chief',
    'member',

    'tech-lead',
    'moderator',
);

CREATE TABLE city_administration (
    user_id    UUID      PRIMARY KEY,
    city_id    UUID      NOT NULL REFERENCES city(id) ON DELETE CASCADE,

    position   VARCHAR(255),
    label      VARCHAR(255),
    role       city_admin_role NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_city_chief
    ON city_administration (city_id)
    WHERE role = 'chief';

CREATE UNIQUE INDEX IF NOT EXISTS uniq_city_tech_lead
    ON city_administration (city_id)
    WHERE role = 'tech-lead';

CREATE TYPE invite_status AS ENUM (
    'sent',
    'accepted'
);

CREATE TABLE invites (
    id           UUID              PRIMARY KEY,
    user_id      UUID              NOT NULL,
    city_id      UUID              NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    initiator_id UUID              NOT NULL,
    status       invite_status     NOT NULL DEFAULT 'sent',
    role         city_admin_role NOT NULL,

    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
);

CREATE OR REPLACE FUNCTION check_city_can_have_admin()
RETURNS trigger AS $$
DECLARE
    st city_status;
BEGIN
    SELECT status INTO st
    FROM cities
    WHERE id = NEW.city_id;

    IF st IN ('unsupported', 'suspended') THEN
        RAISE EXCEPTION
            'City % has status %, admins are not allowed',
            NEW.city_id, st;
    END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_city_administration_ins_upd
    BEFORE INSERT OR UPDATE OF city_id
    ON city_administration
    FOR EACH ROW
    EXECUTE FUNCTION check_city_can_have_admin();

CREATE OR REPLACE FUNCTION check_city_can_have_invite()
RETURNS trigger AS $$
DECLARE
    st city_status;
BEGIN
    SELECT status INTO st
    FROM cities
    WHERE id = NEW.city_id;

    IF st IN ('unsupported', 'suspended') THEN
        RAISE EXCEPTION
            'City % has status %, invites are not allowed',
            NEW.city_id, st;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_invites_ins_upd
    BEFORE INSERT OR UPDATE OF city_id
    ON invites
    FOR EACH ROW
    EXECUTE FUNCTION check_city_can_have_invite();

CREATE OR REPLACE FUNCTION check_city_status_change()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'unsupported' THEN
        IF EXISTS (
            SELECT 1 FROM city_administration a
            WHERE a.city_id = NEW.id
        ) THEN
            RAISE EXCEPTION
                'City % has admins, cannot set status %',
                NEW.id, NEW.status;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cities_status_change
    BEFORE UPDATE OF status
    ON cities
    FOR EACH ROW
    EXECUTE FUNCTION check_city_status_change();

-- +migrate Down
DROP TRIGGER IF EXISTS trg_cities_status_change ON cities;
DROP FUNCTION IF EXISTS check_city_status_change;

DROP TRIGGER IF EXISTS trg_invites_ins_upd ON invites;
DROP FUNCTION IF EXISTS check_city_can_have_invite;

DROP TRIGGER IF EXISTS trg_city_administration_ins_upd ON city_administration;
DROP FUNCTION IF EXISTS check_city_can_have_admin;

DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS city_administration CASCADE;

DROP TYPE IF EXISTS invite_status;
DROP TYPE IF EXISTS city_admin_role;