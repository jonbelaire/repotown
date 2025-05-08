-- Drop tables in reverse order
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS parking_sessions;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS garages;

-- Drop enum types
DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS parking_session_status;
DROP TYPE IF EXISTS vehicle_type;
DROP TYPE IF EXISTS garage_status;