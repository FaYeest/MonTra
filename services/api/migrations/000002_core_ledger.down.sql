-- 000002_core_ledger.down.sql
-- Rolls back core financial ledger schema

DROP TRIGGER IF EXISTS trg_transfers_updated_at ON transfers;
DROP TRIGGER IF EXISTS trg_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS trg_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

DROP TABLE IF EXISTS transfers CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS users CASCADE;
