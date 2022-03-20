CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users
(
    "id"         UUID                  DEFAULT uuid_generate_v4(),
    "login"      VARCHAR(100) NOT NULL,
    "password"   VARCHAR      NOT NULL,
    "created_at" TIMESTAMPTZ  NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ  NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX users_login_idx ON users (login) WHERE deleted_at IS NULL;

-- Orders table
CREATE TABLE orders
(
    "id"          UUID                 DEFAULT uuid_generate_v4(),
    "user_id"     UUID        NOT NULL,
    "number"      VARCHAR(16) NOT NULL,
    "status"      VARCHAR(10) NOT NULL DEFAULT 'NEW',
    "accrual"     INT         NOT NULL,
    "uploaded_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    UNIQUE ("number")
);

CREATE INDEX orders_user_id_idx ON orders ("user_id");

-- Transactions table
CREATE TABLE transactions
(
    "id"           UUID                 DEFAULT uuid_generate_v4(),
    "user_id"      UUID        NOT NULL,
    "order"        VARCHAR(16) NOT NULL,
    "accrual"      INT         NOT NULL,
    "withdrawal"   INT         NOT NULL,
    "processed_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE INDEX transactions_user_id_idx ON transactions ("user_id");
CREATE INDEX transactions_order_idx ON transactions ("order");
CREATE INDEX transactions_processed_at_idx ON transactions (processed_at);
