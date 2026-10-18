# Database Migrations (`services/api/migrations`)

PostgreSQL is the authoritative source of truth for all MonTra financial data.

## Conventions

1. **Sequential Naming:** Each migration consists of an `.up.sql` and `.down.sql` pair prefixed with a 6-digit sequence number (`000001_name.up.sql` / `000001_name.down.sql`).
2. **Monetary Precision:** Never use `REAL`, `FLOAT`, or `DOUBLE PRECISION` for monetary values. Use `BIGINT` (integer units, e.g., whole Rupiah for `IDR`) alongside an explicit `currency VARCHAR(3) NOT NULL` column.
3. **Positive Amount Invariant:** Transaction and transfer principal amounts must have `CHECK (amount > 0)`. Transfer fees must have `CHECK (fee_amount >= 0)`.
4. **Internal Transfer Isolation:** `transfers` must enforce `CHECK (source_account_id <> destination_account_id)` and reference accounts owned by the same user.
5. **Reversibility:** Every `.up.sql` migration must have a tested `.down.sql` counterpart.
