CREATE EXTENSION IF NOT Exists "uuid-ossp";

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE trans_status AS ENUM('pending', 'approved', 'rejected', 'fraud');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_wal_id UUID NOT NULL REFERENCES wallets(id),
    to_wal_id UUID NOT NULL REFERENCES wallets(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    status trans_status NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trans_from_wal ON transactions(from_wal_id);
CREATE INDEX idx_status ON transactions(status);
CREATE INDEX idx_creation_time ON transactions(created_at);