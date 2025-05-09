-- Drop tables first (in reverse order of creation to handle dependencies)
DROP TABLE IF EXISTS migrations;
DROP TABLE IF EXISTS heartbeats;
DROP TABLE IF EXISTS locks;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS operations;
DROP TABLE IF EXISTS intents;

-- Drop ENUMs
DROP TYPE IF EXISTS operation_type;
DROP TYPE IF EXISTS operation_status;
DROP TYPE IF EXISTS intent_status;