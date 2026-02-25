DROP INDEX IF EXISTS idx_creation_time;
DROP INDEX IF EXISTS idx_status;
DROP INDEX IF EXISTS idx_trans_from_wal;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;
DROP TYPE IF EXISTS trans_status;
DROP EXTENSION IF EXISTS "uuid-ossp";