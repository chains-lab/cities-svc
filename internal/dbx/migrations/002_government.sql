-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_gov_roles AS ENUM (
  'mayor',
  'advisor',
  'member',
  'moderator'
);

CREATE TYPE city_gov_statuses AS ENUM ('active', 'inactive');

CREATE TABLE city_governments (
    id             UUID              PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID              NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    city_id        UUID              NOT NULL REFERENCES city(id)   ON DELETE CASCADE,
    status         city_gov_statuses NOT NULL,
    role           city_gov_roles    NOT NULL,
    label          VARCHAR(255)      NOT NULL,
    deactivated_at TIMESTAMP,
    created_at     TIMESTAMP         NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at     TIMESTAMP         NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    CONSTRAINT city_gov_status_deactivated_ck CHECK (
        (status = 'active'   AND deactivated_at IS NULL) OR
        (status = 'inactive' AND deactivated_at IS NOT NULL)
    )
);

CREATE UNIQUE INDEX city_gov_unique_mayor_active
    ON city_governments(city_id)
    WHERE role = 'mayor' AND status = 'active';

CREATE UNIQUE INDEX city_gov_user_unique_global_active
    ON city_governments(user_id)
    WHERE status = 'active';

CREATE INDEX city_gov_city_status_role_idx
    ON city_governments(city_id, status, role);

-- +migrate Down
DROP INDEX IF EXISTS city_gov_city_status_role_idx;
DROP INDEX IF EXISTS city_gov_user_unique_global_active;
DROP INDEX IF EXISTS city_gov_unique_mayor_active;

DROP TABLE IF EXISTS city_governments CASCADE;
DROP TYPE IF EXISTS city_gov_roles;
DROP TYPE IF EXISTS city_gov_statuses;
DROP EXTENSION IF EXISTS "uuid-ossp";
