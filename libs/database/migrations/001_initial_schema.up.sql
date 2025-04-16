-- IntentSchema
CREATE TABLE IF NOT EXISTS intents (
    id BIGSERIAL PRIMARY KEY,
    signature TEXT,
    identity TEXT,
    identity_curve TEXT,
    status TEXT,
    exipry BIGINT,
    created_at BIGINT
);

-- OperationSchema
CREATE TABLE IF NOT EXISTS operations (
    id BIGSERIAL PRIMARY KEY,
    intent_id BIGINT REFERENCES intents(id),
    serialized_txn TEXT,
    data_to_sign TEXT,
    chain_id TEXT,
    genesis_hash TEXT,
    key_curve TEXT,
    status TEXT,
    result TEXT,
    type TEXT,
    solver TEXT,
    solver_metadata TEXT,
    solver_data_to_sign TEXT,
    solver_output TEXT
);

-- WalletSchema
CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    identity TEXT,
    identity_curve TEXT,
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
    signers TEXT
);

-- LockSchema
CREATE TABLE IF NOT EXISTS locks (
    id BIGSERIAL PRIMARY KEY,
    identity TEXT,
    identity_curve TEXT,
    locked BOOLEAN
);

-- HeartbeatSchema
CREATE TABLE IF NOT EXISTS heartbeats (
    publickey TEXT PRIMARY KEY,
    "timestamp" TIMESTAMPTZ
);

-- Migrations Schema
CREATE TABLE IF NOT EXISTS migrations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
