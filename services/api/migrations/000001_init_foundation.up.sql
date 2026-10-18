-- 000001_init_foundation.up.sql
-- Initial foundation migration for MonTra PostgreSQL database.
-- Enables UUID generation and establishes schema metadata tracking.
-- Complete financial domain tables (users, accounts, categories, transactions, transfers, budgets)
-- are scheduled for Phase 1 & Phase 2 migrations per docs/IMPLEMENTATION_PLAN.md.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
