BEGIN;

CREATE TABLE accounts (
  id VARCHAR NOT NULL,
  owner VARCHAR NOT NULL,
  balance VARCHAR(64) NOT NULL DEFAULT '0',
  currency VARCHAR(16) NOT NULL
);

CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  account VARCHAR NOT NULL,
  from_account VARCHAR,
  to_account VARCHAR,
  direction VARCHAR(16) NOT NULL,
  amount VARCHAR(64) NOT NULL
);

CREATE UNIQUE INDEX ON accounts(id);

COMMIT;