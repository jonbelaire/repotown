# Banking Microservice

A microservice for managing fictitious banking operations in Repotown.

## Architecture

The banking microservice follows a layered architecture:

1. **API Layer**: RESTful API endpoints for client interactions
2. **Service Layer**: Business logic and transaction handling
3. **Repository Layer**: Data access and persistence
4. **Domain Layer**: Core domain models and business rules

### Component Structure

```
services/banking/
├── cmd/                  # Application entry points
│   └── server/           # HTTP server
├── internal/             # Private application code
│   ├── api/              # API handlers and middleware
│   ├── domain/           # Core domain models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── config/           # Service configuration
├── migrations/           # Database migrations
├── scripts/              # Build and deployment scripts
└── pkg/                  # Public packages (if any)
```

## Core Domain Entities

- **Account**: Represents a bank account with balance and transaction history
- **Transaction**: Records money movements between accounts
- **Customer**: Account owner information
- **Card**: Payment card associated with accounts

## API Endpoints

### Accounts
- `GET /api/accounts` - List accounts
- `GET /api/accounts/{id}` - Get account details
- `POST /api/accounts` - Create new account
- `PUT /api/accounts/{id}` - Update account
- `DELETE /api/accounts/{id}` - Close account

### Transactions
- `GET /api/transactions` - List transactions (with filtering)
- `GET /api/transactions/{id}` - Get transaction details
- `POST /api/transactions` - Create transaction (transfer, deposit, withdrawal)

### Customers
- `GET /api/customers` - List customers
- `GET /api/customers/{id}` - Get customer details
- `POST /api/customers` - Create customer
- `PUT /api/customers/{id}` - Update customer
- `DELETE /api/customers/{id}` - Delete customer

## Technology Stack

- **Language**: Go 1.20+
- **API Framework**: Chi router with go-core utilities
- **Database**: PostgreSQL
- **Authentication**: JWT-based auth using go-core/auth
- **Configuration**: Environment-based using go-core/config
- **Logging**: Structured logging with go-core/logging

## Development

```bash
# Start the service
go run cmd/server/main.go

# Run tests
go test ./...
```