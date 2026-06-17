CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL CHECK (balance >= 0)
);

CREATE TABLE transfers (
    id UUID PRIMARY KEY,
    idempotency_key TEXT UNIQUE NOT NULL,
    from_wallet_id UUID NOT NULL,
    to_wallet_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    transfer_id UUID NOT NULL REFERENCES transfers(id),
    wallet_id UUID NOT NULL,
    type TEXT NOT NULL,
    amount BIGINT NOT NULL
);

CREATE TABLE idempotency_records (
    idempotency_key TEXT PRIMARY KEY,
    transfer_id UUID NOT NULL,
    response JSONB NOT NULL
);

CREATE INDEX idx_wallets_id ON wallets(id);
CREATE INDEX idx_transfers_wallet ON transfers(from_wallet_id, to_wallet_id);
CREATE INDEX idx_ledger_transfer ON ledger_entries(transfer_id);