# Parking Garage Service

This service manages parking garages, vehicles, parking sessions, and reservations.

## Features

- Manage multiple parking garages with real-time space availability
- Track vehicles and their parking sessions
- Support for parking reservations
- Detailed reporting on garage usage

## API Endpoints

### Garages

- `GET /api/garages` - List all garages
- `POST /api/garages` - Create a new garage
- `GET /api/garages/available` - List garages with available spaces
- `GET /api/garages/{id}` - Get garage details
- `PUT /api/garages/{id}` - Update garage details
- `PATCH /api/garages/{id}/status` - Update garage status

### Vehicles

- `GET /api/vehicles` - List all vehicles
- `POST /api/vehicles` - Register a new vehicle
- `GET /api/vehicles/license/{plate}` - Get vehicle by license plate
- `GET /api/vehicles/owner/{ownerID}` - List vehicles by owner
- `GET /api/vehicles/{id}` - Get vehicle details
- `PUT /api/vehicles/{id}` - Update vehicle details
- `PATCH /api/vehicles/{id}/type` - Update vehicle type
- `PATCH /api/vehicles/{id}/owner` - Change vehicle owner

### Parking Sessions

- `GET /api/sessions` - List all parking sessions
- `POST /api/sessions` - Start a new parking session
- `GET /api/sessions/active` - List active parking sessions
- `GET /api/sessions/garage/{garageID}` - List sessions by garage
- `GET /api/sessions/vehicle/{vehicleID}` - List sessions by vehicle
- `GET /api/sessions/{id}` - Get session details
- `PUT /api/sessions/{id}/end` - End a parking session
- `PUT /api/sessions/{id}/cancel` - Cancel a parking session
- `PATCH /api/sessions/{id}/spot` - Update parking spot
- `PATCH /api/sessions/{id}/notes` - Add notes to a session
- `PATCH /api/sessions/{id}/pay` - Mark session as paid

### Reservations

- `GET /api/reservations` - List all reservations
- `POST /api/reservations` - Create a new reservation
- `GET /api/reservations/active` - List active reservations
- `GET /api/reservations/garage/{garageID}` - List reservations by garage
- `GET /api/reservations/vehicle/{vehicleID}` - List reservations by vehicle
- `GET /api/reservations/code/{code}` - Get reservation by confirmation code
- `GET /api/reservations/{id}` - Get reservation details
- `PUT /api/reservations/{id}/use` - Use a reservation
- `PUT /api/reservations/{id}/cancel` - Cancel a reservation
- `POST /api/reservations/check-expired` - Check for expired reservations

## Setup

### Prerequisites

- Go 1.18 or higher
- PostgreSQL 13 or higher
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

### Environment Variables

Create a `.env` file with the following variables:

```
SERVER_ADDRESS=:8082
ENVIRONMENT=development
LOG_LEVEL=info
DATABASE_URL=postgres://username:password@localhost:5432/parking_db?sslmode=disable
JWT_SECRET=your_jwt_secret
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