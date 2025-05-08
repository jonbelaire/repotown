-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_number TEXT,
    street TEXT,
    city TEXT,
    state TEXT,
    country TEXT,
    postal_code TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    currency_code TEXT NOT NULL,
    name TEXT NOT NULL,
    number TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    closed_at TIMESTAMP WITH TIME ZONE
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts(id),
    source_account_id UUID REFERENCES accounts(id),
    target_account_id UUID REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    currency_code TEXT NOT NULL,
    description TEXT NOT NULL,
    reference TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX idx_accounts_customer_id ON accounts(customer_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_source_account_id ON transactions(source_account_id) WHERE source_account_id IS NOT NULL;
CREATE INDEX idx_transactions_target_account_id ON transactions(target_account_id) WHERE target_account_id IS NOT NULL;