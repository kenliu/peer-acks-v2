-- table schema for acks DB
USE peeracks;
DROP TABLE upvotes;
DROP TABLE acks;
DROP TABLE users;

CREATE TABLE users
(
    id              UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    slack_id        STRING UNIQUE,
    slack_name      STRING,
    slack_real_name STRING,
    email           STRING    NOT NULL UNIQUE,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL
);

CREATE TABLE acks
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    msg               STRING    NOT NULL,
    user_id           UUID      NOT NULL REFERENCES users (id),
    source            STRING    NOT NULL,
    slack_msg_ts      STRING,
    slack_msg_channel STRING,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL
);

CREATE TABLE upvotes
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reaction       STRING NOT NULL, -- the reactji name
    ack_id         UUID   NOT NULL REFERENCES acks (id),
    user_id        UUID   NOT NULL REFERENCES users (id),
    slack_event_ts STRING NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);