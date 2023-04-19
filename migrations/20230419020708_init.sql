-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE file_status AS ENUM ('created', 'uploading', 'uploaded');

CREATE TABLE IF NOT EXISTS files
(
    id           uuid PRIMARY KEY UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    name         varchar(255)            NOT NULL,
    content_type varchar(255)            NOT NULL,
    is_uploaded  bool                    NOT NULL DEFAULT false,
    created_at   timestamptz             NOT NULL DEFAULT NOW(),
    updated_at   timestamptz             NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chunks
(
    id         uuid PRIMARY KEY UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    file_id    uuid                    NOT NULL REFERENCES files,
    number     smallint                NOT NULL,
    size       bigint                  NOT NULL,
    check_sum  varchar(255)            NOT NULL,
    url        varchar(255)            NOT NULL,
    bucket     varchar(255)            NOT NULL,
    created_at timestamptz             NOT NULL DEFAULT NOW(),
    updated_at timestamptz             NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS chunks;
DROP TABLE IF EXISTS files;
DROP TYPE file_status;
