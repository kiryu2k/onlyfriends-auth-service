-- +goose Up
CREATE TABLE auth_user
(
    id              TEXT PRIMARY KEY,
    email           TEXT UNIQUE                 NOT NULL,
    hashed_password TEXT                        NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE auth_user;
