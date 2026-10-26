-- 000002_core_ledger.up.sql
-- Core Financial Ledger Schema for MonTra
-- Implements tables: users, accounts, categories, transactions, transfers
-- Strictly enforces exact integer BIGINT for currency, precise VARCHAR lengths, and referential constraints.

-- 1. Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash CHAR(60) NOT NULL,
    display_name VARCHAR(60) NOT NULL,
    preferred_currency CHAR(3) NOT NULL DEFAULT 'IDR',
    timezone VARCHAR(32) NOT NULL DEFAULT 'Asia/Jakarta',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- 2. Accounts Table (Wallets, Bank Accounts, E-Wallets)
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    account_type VARCHAR(10) NOT NULL CHECK (account_type IN ('cash', 'bank', 'ewallet', 'other')),
    opening_balance BIGINT NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_user_archived ON accounts(user_id, archived_at) WHERE archived_at IS NULL;

CREATE TRIGGER trg_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- 3. Categories Table (System Default and User Custom)
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- NULL indicates system-provided default category
    name VARCHAR(50) NOT NULL,
    transaction_type VARCHAR(7) NOT NULL CHECK (transaction_type IN ('income', 'expense')),
    icon VARCHAR(32),
    color VARCHAR(9),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_user_type ON categories(user_id, transaction_type);
CREATE INDEX idx_categories_system ON categories(is_system, transaction_type);

CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- 4. Transactions Table (Income and Expense)
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
    type VARCHAR(7) NOT NULL CHECK (type IN ('income', 'expense')),
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    description VARCHAR(120) NOT NULL,
    transaction_date TIMESTAMPTZ NOT NULL,
    notes TEXT,
    source VARCHAR(8) NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'import', 'ocr')),
    category_source VARCHAR(8) NOT NULL DEFAULT 'manual' CHECK (category_source IN ('manual', 'model', 'system')),
    category_model_version VARCHAR(32),
    idempotency_key VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE UNIQUE INDEX idx_transactions_user_idempotency ON transactions(user_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

CREATE TRIGGER trg_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- 5. Transfers Table (Atomic Internal Fund Movements)
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    destination_account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    fee_amount BIGINT NOT NULL DEFAULT 0 CHECK (fee_amount >= 0),
    transaction_date TIMESTAMPTZ NOT NULL,
    description VARCHAR(120),
    idempotency_key VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_transfers_distinct_accounts CHECK (source_account_id <> destination_account_id)
);

CREATE INDEX idx_transfers_user_date ON transfers(user_id, transaction_date DESC);
CREATE INDEX idx_transfers_source ON transfers(source_account_id);
CREATE INDEX idx_transfers_destination ON transfers(destination_account_id);
CREATE UNIQUE INDEX idx_transfers_user_idempotency ON transfers(user_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

CREATE TRIGGER trg_transfers_updated_at
    BEFORE UPDATE ON transfers
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- 6. Seed Default System Categories (is_system = true, user_id = NULL)
-- Expense categories
INSERT INTO categories (name, transaction_type, icon, color, is_system) VALUES
('Makanan & Minuman', 'expense', 'utensils', '#FF7043', TRUE),
('Transportasi', 'expense', 'car', '#42A5F5', TRUE),
('Belanja Harian', 'expense', 'shopping-bag', '#26A69A', TRUE),
('Tagihan & Utilitas', 'expense', 'receipt', '#AB47BC', TRUE),
('Kesehatan', 'expense', 'heart-pulse', '#EF5350', TRUE),
('Pendidikan', 'expense', 'graduation-cap', '#5C6BC0', TRUE),
('Hiburan', 'expense', 'film', '#FFA726', TRUE),
('Donasi & Sosial', 'expense', 'hand-heart', '#8D6E63', TRUE),
('Pengeluaran Lainnya', 'expense', 'dots-three', '#78909C', TRUE);

-- Income categories
INSERT INTO categories (name, transaction_type, icon, color, is_system) VALUES
('Gaji', 'income', 'briefcase', '#66BB6A', TRUE),
('Uang Saku', 'income', 'wallet', '#26C6DA', TRUE),
('Freelance', 'income', 'laptop', '#7E57C2', TRUE),
('Bisnis', 'income', 'storefront', '#FFA000', TRUE),
('Hadiah & Bonus', 'income', 'gift', '#EC407A', TRUE),
('Investasi', 'income', 'trend-up', '#2E7D32', TRUE),
('Pemasukan Lainnya', 'income', 'plus-circle', '#546E7A', TRUE);
