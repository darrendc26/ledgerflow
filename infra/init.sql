CREATE SEQUENCE account_number_seq START 10000000;

CREATE TABLE accounts (
    id TEXT PRIMARY KEY DEFAULT nextval('account_number_seq')::text,
    user_id TEXT NOT NULL,
    balance BIGINT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    reference_id TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE payments (
    id TEXT PRIMARY KEY,
    sender_account TEXT NOT NULL REFERENCES accounts(id),
    receiver_account TEXT NOT NULL REFERENCES accounts(id),
    transaction_id TEXT REFERENCES transactions(id),
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE SEQUENCE ledger_entries_id_seq START 1;
CREATE TABLE ledger_entries (
    id TEXT PRIMARY KEY DEFAULT nextval('ledger_entries_id_seq')::text,
    transaction_id TEXT NOT NULL REFERENCES transactions(id),
    account_id TEXT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);