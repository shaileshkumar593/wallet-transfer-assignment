CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL CHECK(balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE transfer_status AS ENUM (
    'PENDING',
    'PROCESSED',
    'FAILED'
);

CREATE TABLE transfers (
    id UUID PRIMARY KEY,

    idempotency_key VARCHAR(255)
    UNIQUE NOT NULL,

    from_wallet_id UUID NOT NULL
    REFERENCES wallets(id),

    to_wallet_id UUID NOT NULL
    REFERENCES wallets(id),

    amount BIGINT NOT NULL
    CHECK(amount > 0),

    status transfer_status NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE ledger_type AS ENUM (
    'DEBIT',
    'CREDIT'
);

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,

    transfer_id UUID NOT NULL
    REFERENCES transfers(id),

    wallet_id UUID NOT NULL
    REFERENCES wallets(id),

    entry_type ledger_type NOT NULL,

    amount BIGINT NOT NULL
    CHECK(amount > 0),

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE idempotency_records (

    idempotency_key VARCHAR(255)
    PRIMARY KEY,

    transfer_id UUID NOT NULL,

    response JSONB NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wallet_balance
ON wallets(balance);

CREATE INDEX idx_transfer_wallets
ON transfers(
    from_wallet_id,
    to_wallet_id
);

CREATE INDEX idx_ledger_transfer
ON ledger_entries(
    transfer_id
);

CREATE UNIQUE INDEX idx_idempotency_key
ON idempotency_records(idempotency_key);