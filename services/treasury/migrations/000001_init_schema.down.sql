-- Drop tables in reverse order
DROP TABLE IF EXISTS tax_payments;
DROP TABLE IF EXISTS filing_credits;
DROP TABLE IF EXISTS filing_deductions;
DROP TABLE IF EXISTS tax_filings;
DROP TABLE IF EXISTS taxpayer_exemptions;
DROP TABLE IF EXISTS taxpayers;
DROP TABLE IF EXISTS tax_rates;

-- Drop enum types
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS filing_period;
DROP TYPE IF EXISTS filing_status;
DROP TYPE IF EXISTS taxpayer_status;
DROP TYPE IF EXISTS taxpayer_type;
DROP TYPE IF EXISTS tax_bracket_type;
DROP TYPE IF EXISTS tax_status;
DROP TYPE IF EXISTS tax_type;