# Treasury Service

This service manages tax revenue and treasury matters for the fictional town, including setting tax rates, collecting income taxes, and processing sales taxes from all businesses.

## Features

- Tax rate management (income, sales, property, business, excise)
- Taxpayer registration and management
- Tax filing submission and processing
- Tax payment processing
- Comprehensive reporting and analytics

## API Endpoints

### Tax Rates

- `GET /api/tax-rates` - List all tax rates
- `POST /api/tax-rates` - Create a new tax rate
- `GET /api/tax-rates/active` - List active tax rates
- `GET /api/tax-rates/type/{type}` - List tax rates by type
- `GET /api/tax-rates/jurisdiction/{jurisdictionCode}` - List tax rates by jurisdiction
- `GET /api/tax-rates/calculate/income` - Calculate income tax
- `GET /api/tax-rates/calculate/sales` - Calculate sales tax
- `GET /api/tax-rates/{id}` - Get tax rate details
- `PUT /api/tax-rates/{id}` - Update tax rate
- `POST /api/tax-rates/{id}/activate` - Activate tax rate
- `POST /api/tax-rates/{id}/deactivate` - Deactivate tax rate
- `POST /api/tax-rates/{id}/archive` - Archive tax rate

### Taxpayers

- `GET /api/taxpayers` - List all taxpayers
- `POST /api/taxpayers` - Register a new taxpayer
- `GET /api/taxpayers/search` - Search for taxpayers
- `GET /api/taxpayers/type/{type}` - List taxpayers by type
- `GET /api/taxpayers/status/{status}` - List taxpayers by status
- `GET /api/taxpayers/industry/{industry}` - List businesses by industry
- `GET /api/taxpayers/identifier/{identifier}` - Get taxpayer by tax identifier
- `GET /api/taxpayers/{id}` - Get taxpayer details
- `PATCH /api/taxpayers/{id}/status` - Update taxpayer status
- `PATCH /api/taxpayers/{id}/contact` - Update taxpayer contact info
- `PATCH /api/taxpayers/{id}/address` - Update taxpayer address
- `PATCH /api/taxpayers/{id}/business` - Update business info
- `POST /api/taxpayers/{id}/exemptions` - Add tax exemption
- `DELETE /api/taxpayers/{id}/exemptions/{code}` - Remove tax exemption

### Tax Filings

- `GET /api/filings` - List all tax filings
- `POST /api/filings` - Create a new tax filing
- `GET /api/filings/status/{status}` - List filings by status
- `GET /api/filings/taxpayer/{taxpayerID}` - List filings by taxpayer
- `GET /api/filings/period/{year}/{period}` - List filings by period
- `GET /api/filings/overdue` - List overdue filings
- `GET /api/filings/recent` - List recently submitted filings
- `GET /api/filings/{id}` - Get filing details
- `PATCH /api/filings/{id}/amounts` - Update filing amounts
- `POST /api/filings/{id}/deductions` - Add deduction
- `POST /api/filings/{id}/credits` - Add credit
- `POST /api/filings/{id}/submit` - Submit filing
- `POST /api/filings/{id}/process` - Process filing
- `POST /api/filings/{id}/accept` - Accept filing
- `POST /api/filings/{id}/reject` - Reject filing
- `POST /api/filings/{id}/amend` - Amend filing

### Tax Payments

- `GET /api/payments` - List all tax payments
- `POST /api/payments` - Create a new tax payment
- `GET /api/payments/code/{code}` - Get payment by confirmation code
- `GET /api/payments/taxpayer/{taxpayerID}` - List payments by taxpayer
- `GET /api/payments/filing/{filingID}` - List payments for a filing
- `GET /api/payments/status/{status}` - List payments by status
- `GET /api/payments/date-range` - List payments within a date range
- `GET /api/payments/recent` - List recent payments
- `GET /api/payments/totals/by-type` - Get total payments by tax type
- `GET /api/payments/{id}` - Get payment details
- `POST /api/payments/{id}/process` - Process payment
- `POST /api/payments/{id}/fail` - Mark payment as failed
- `POST /api/payments/{id}/refund` - Refund payment
- `POST /api/payments/{id}/void` - Void payment
- `PATCH /api/payments/{id}/amount` - Update payment amount

### Reports

- `GET /api/reports/revenue` - Generate revenue report
- `GET /api/reports/filing-status` - Generate filing status report
- `GET /api/reports/taxpayer-compliance` - Generate taxpayer compliance report
- `GET /api/reports/tax-type-breakdown` - Generate tax type breakdown report

## Setup

### Prerequisites

- Go 1.18 or higher
- PostgreSQL 13 or higher
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

### Environment Variables

Create a `.env` file with the following variables:

```
SERVER_ADDRESS=:8083
ENVIRONMENT=development
LOG_LEVEL=info
DATABASE_URL=postgres://username:password@localhost:5432/treasury_db?sslmode=disable
JWT_SECRET=your_jwt_secret
DEFAULT_JURISDICTION=CITYTOWN
```

### Running the Service

1. Apply database migrations:
   ```
   make migrate-up
   ```

2. Build and run the service:
   ```
   make build
   make run
   ```

### Docker

You can also run the service using Docker:

```
make docker-build
make docker-run
```

## Testing

Run the test suite with:

```
make test
```