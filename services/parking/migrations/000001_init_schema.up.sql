-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE garage_status AS ENUM ('operational', 'maintenance', 'closed');
CREATE TYPE vehicle_type AS ENUM ('car', 'motorcycle', 'truck', 'electric');
CREATE TYPE parking_session_status AS ENUM ('active', 'completed', 'cancelled');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');
CREATE TYPE reservation_status AS ENUM ('active', 'used', 'cancelled', 'expired');

-- Create garages table
CREATE TABLE garages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    status garage_status NOT NULL DEFAULT 'operational',
    total_spaces INTEGER NOT NULL CHECK (total_spaces > 0),
    available_spaces INTEGER NOT NULL CHECK (available_spaces >= 0),
    hourly_rate BIGINT NOT NULL CHECK (hourly_rate >= 0),
    daily_rate BIGINT NOT NULL CHECK (daily_rate >= 0),
    operating_hours VARCHAR(255) NOT NULL,
    has_electric_charging BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_spaces CHECK (available_spaces <= total_spaces)
);

-- Create vehicles table
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    license_plate VARCHAR(50) NOT NULL UNIQUE,
    type vehicle_type NOT NULL,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    color VARCHAR(50) NOT NULL,
    owner_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create parking_sessions table
CREATE TABLE parking_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    garage_id UUID NOT NULL REFERENCES garages(id),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    spot_number VARCHAR(50),
    status parking_session_status NOT NULL DEFAULT 'active',
    payment_status payment_status NOT NULL DEFAULT 'pending',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    duration INTERVAL,
    amount_charged BIGINT,
    is_prepaid BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create reservations table
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    garage_id UUID NOT NULL REFERENCES garages(id),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    status reservation_status NOT NULL DEFAULT 'active',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    amount_paid BIGINT NOT NULL CHECK (amount_paid >= 0),
    confirmation_code VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    CONSTRAINT check_timespan CHECK (end_time > start_time)
);

-- Create indexes
CREATE INDEX idx_parking_sessions_garage_id ON parking_sessions(garage_id);
CREATE INDEX idx_parking_sessions_vehicle_id ON parking_sessions(vehicle_id);
CREATE INDEX idx_parking_sessions_status ON parking_sessions(status);
CREATE INDEX idx_reservations_garage_id ON reservations(garage_id);
CREATE INDEX idx_reservations_vehicle_id ON reservations(vehicle_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_confirmation_code ON reservations(confirmation_code);
CREATE INDEX idx_vehicles_license_plate ON vehicles(license_plate);
CREATE INDEX idx_vehicles_owner_id ON vehicles(owner_id);
CREATE INDEX idx_garages_status ON garages(status);