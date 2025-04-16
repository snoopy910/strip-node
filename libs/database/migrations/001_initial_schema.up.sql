-- Create ENUM types first
CREATE TYPE intent_status AS ENUM (
    'PROCESSING',
    'COMPLETED',
    'FAILED',
    'EXPIRED'
);

CREATE TYPE operation_status AS ENUM (
    'PENDING',
    'WAITING',
    'COMPLETED',
    'FAILED',
    'EXPIRED'
);

CREATE TYPE operation_type AS ENUM (
    'TRANSACTION',
    'SOLVER',
    'BRIDGE_DEPOSIT',
    'SWAP',
    'BURN',
    'BURN_SYNTHETIC',
    'WITHDRAW',
    'SEND_TO_BRIDGE'
);

-- IntentSchema
CREATE TABLE IF NOT EXISTS intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    signature TEXT NOT NULL,
    identity TEXT NOT NULL,
    identity_curve TEXT NOT NULL,
    status intent_status NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- OperationSchema
CREATE TABLE IF NOT EXISTS operations (
    id BIGSERIAL PRIMARY KEY,
    intent_id UUID REFERENCES intents(id) NOT NULL,
    serialized_txn TEXT NOT NULL,
    data_to_sign TEXT NOT NULL,
    chain_id TEXT NOT NULL,
    genesis_hash TEXT,
    key_curve TEXT NOT NULL,
    status operation_status NOT NULL,
    result TEXT,
    type operation_type NOT NULL,
    solver TEXT,
    solver_metadata JSONB,
    solver_data_to_sign TEXT,
    solver_output JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- WalletSchema
CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    identity_curve TEXT NOT NULL,
    eddsa_public_key TEXT,
    aptos_eddsa_public_key TEXT,
    ecdsa_public_key TEXT,
    bitcoin_mainnet_public_key TEXT,
    bitcoin_testnet_public_key TEXT,
    bitcoin_regtest_public_key TEXT,
    stellar_public_key TEXT,
    dogecoin_mainnet_public_key TEXT,
    dogecoin_testnet_public_key TEXT,
    sui_public_key TEXT,
    algorand_eddsa_public_key TEXT,
    ripple_eddsa_public_key TEXT,
    cardano_public_key TEXT,
    signers JSONB
);

-- LockSchema
CREATE TABLE IF NOT EXISTS locks (
    id BIGSERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    identity_curve TEXT NOT NULL,
    locked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_locks_identity ON locks(identity);

-- HeartbeatSchema
CREATE TABLE IF NOT EXISTS heartbeats (
    publickey TEXT PRIMARY KEY,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Migrations Schema
CREATE TABLE IF NOT EXISTS migrations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
