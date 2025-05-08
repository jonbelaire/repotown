-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE tax_type AS ENUM ('income', 'sales', 'property', 'business', 'excise');
CREATE TYPE tax_status AS ENUM ('active', 'inactive', 'proposed', 'archived');
CREATE TYPE tax_bracket_type AS ENUM ('flat', 'progressive', 'tiered');
CREATE TYPE taxpayer_type AS ENUM ('individual', 'business', 'non_profit', 'government');
CREATE TYPE taxpayer_status AS ENUM ('active', 'inactive', 'exempt', 'delinquent');
CREATE TYPE filing_status AS ENUM ('draft', 'submitted', 'processing', 'accepted', 'rejected', 'amended', 'audited');
CREATE TYPE filing_period AS ENUM ('monthly', 'quarterly', 'semi_annual', 'annual');
CREATE TYPE payment_method AS ENUM ('electronic', 'check', 'credit', 'debit', 'cash', 'wire');
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'refunded', 'voided');

-- Create tax_rates table
CREATE TABLE tax_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type tax_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    rate DECIMAL(5, 4) NOT NULL CHECK (rate >= 0 AND rate <= 1),
    bracket_type tax_bracket_type NOT NULL,
    min_amount BIGINT,
    max_amount BIGINT,
    status tax_status NOT NULL DEFAULT 'proposed',
    category VARCHAR(100),
    jurisdiction_code VARCHAR(50) NOT NULL,
    effective_date TIMESTAMPTZ NOT NULL,
    expiration_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_amount_range CHECK (min_amount IS NULL OR max_amount IS NULL OR max_amount > min_amount)
);

-- Create taxpayers table
CREATE TABLE taxpayers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type taxpayer_type NOT NULL,
    status taxpayer_status NOT NULL DEFAULT 'active',
    name VARCHAR(255) NOT NULL,
    tax_identifier VARCHAR(100) NOT NULL UNIQUE,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50) NOT NULL,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    annual_revenue BIGINT,
    business_type VARCHAR(100),
    industry VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create taxpayer_exemptions table
CREATE TABLE taxpayer_exemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id),
    exemption_code VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(taxpayer_id, exemption_code)
);

-- Create tax_filings table
CREATE TABLE tax_filings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id),
    tax_year INTEGER NOT NULL CHECK (tax_year >= 1900 AND tax_year <= 2100),
    period filing_period NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    filing_type tax_type NOT NULL,
    status filing_status NOT NULL DEFAULT 'draft',
    gross_income BIGINT,
    taxable_income BIGINT,
    total_sales BIGINT,
    taxable_amount BIGINT,
    tax_calculated BIGINT DEFAULT 0,
    tax_paid BIGINT DEFAULT 0,
    submission_date TIMESTAMPTZ,
    acceptance_date TIMESTAMPTZ,
    due_date TIMESTAMPTZ NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_period CHECK (period_end > period_start)
);

-- Create filing_deductions table
CREATE TABLE filing_deductions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filing_id UUID NOT NULL REFERENCES tax_filings(id),
    code VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create filing_credits table
CREATE TABLE filing_credits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filing_id UUID NOT NULL REFERENCES tax_filings(id),
    code VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create tax_payments table
CREATE TABLE tax_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id),
    filing_id UUID REFERENCES tax_filings(id),
    tax_type tax_type NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    payment_method payment_method NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    payment_date TIMESTAMPTZ NOT NULL,
    confirmation_code VARCHAR(50) NOT NULL UNIQUE,
    notes TEXT,
    processed_at TIMESTAMPTZ,
    refunded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_tax_rates_type ON tax_rates(type);
CREATE INDEX idx_tax_rates_jurisdiction ON tax_rates(jurisdiction_code);
CREATE INDEX idx_tax_rates_status ON tax_rates(status);
CREATE INDEX idx_tax_rates_effective_date ON tax_rates(effective_date);

CREATE INDEX idx_taxpayers_type ON taxpayers(type);
CREATE INDEX idx_taxpayers_status ON taxpayers(status);
CREATE INDEX idx_taxpayers_tax_identifier ON taxpayers(tax_identifier);
CREATE INDEX idx_taxpayers_industry ON taxpayers(industry);

CREATE INDEX idx_tax_filings_taxpayer ON tax_filings(taxpayer_id);
CREATE INDEX idx_tax_filings_year_period ON tax_filings(tax_year, period);
CREATE INDEX idx_tax_filings_status ON tax_filings(status);
CREATE INDEX idx_tax_filings_due_date ON tax_filings(due_date);
CREATE INDEX idx_tax_filings_submission_date ON tax_filings(submission_date);

CREATE INDEX idx_filing_deductions_filing ON filing_deductions(filing_id);
CREATE INDEX idx_filing_credits_filing ON filing_credits(filing_id);

CREATE INDEX idx_tax_payments_taxpayer ON tax_payments(taxpayer_id);
CREATE INDEX idx_tax_payments_filing ON tax_payments(filing_id);
CREATE INDEX idx_tax_payments_status ON tax_payments(status);
CREATE INDEX idx_tax_payments_payment_date ON tax_payments(payment_date);
CREATE INDEX idx_tax_payments_confirmation_code ON tax_payments(confirmation_code);