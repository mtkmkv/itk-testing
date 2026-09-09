CREATE SCHEMA IF NOT EXISTS wallet;

-- ============================================================
-- 1. КОШЕЛЬКИ
-- ============================================================

CREATE TABLE IF NOT EXISTS wallet.wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT wallets_balance_non_negative
        CHECK (balance >= 0)
);

-- ============================================================
-- 2. ОПЕРАЦИИ КОШЕЛЬКА
-- ============================================================

CREATE TABLE IF NOT EXISTS wallet.operations (
    id BIGSERIAL PRIMARY KEY,

    wallet_id UUID NOT NULL
        REFERENCES wallet.wallets(id)
        ON DELETE CASCADE,

    operation_type VARCHAR(10) NOT NULL,

    amount BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT operations_type_check
        CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),

    CONSTRAINT operations_amount_check
        CHECK (amount > 0)
);

-- ============================================================
-- 3. ИНДЕКСЫ
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_operations_wallet_id
    ON wallet.operations(wallet_id);

CREATE INDEX IF NOT EXISTS idx_operations_wallet_id_created_at
    ON wallet.operations(wallet_id, created_at);