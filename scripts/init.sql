-- Pulse Patrol - Aurora PostgreSQL DDL Schema
-- Date: 2026-02-17
-- Target: Amazon Aurora (PostgreSQL) compatible with Docker Compose

-- Enable UUID extension for secure, non-sequential identifiers
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. INSTITUTIONAL & ACCESS CONTROL (FR75)
-- ==========================================

-- Stores medical institution profiles (Tenant Isolation)
CREATE TABLE institutions (
                              institution_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              name VARCHAR(255) NOT NULL,
                              license_key VARCHAR(100) UNIQUE NOT NULL,
                              security_policy_json JSONB, -- Stores institutional security metadata (NFR36)
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Defines the available roles (Extensible RBAC)
CREATE TABLE roles (
                       role_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       role_name VARCHAR(50) UNIQUE NOT NULL, -- e.g., 'ADMIN', 'DOCTOR', 'NURSE'
                       description TEXT,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Defines granular actions (FR77 Mapping)
CREATE TABLE permissions (
                             permission_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             permission_key VARCHAR(100) UNIQUE NOT NULL, -- e.g., 'READ_TELEMETRY'
                             description TEXT,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Mapping Table (Many-to-Many)
CREATE TABLE role_permissions (
                                  role_id UUID REFERENCES roles(role_id) ON DELETE CASCADE,
                                  permission_id UUID REFERENCES permissions(permission_id) ON DELETE CASCADE,
                                  PRIMARY KEY (role_id, permission_id)
);

-- Medical staff accounts
CREATE TABLE staff_users (
                             user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             institution_id UUID REFERENCES institutions(institution_id) ON DELETE CASCADE,
                             role_id UUID REFERENCES roles(role_id) ON DELETE RESTRICT,
                             email VARCHAR(255) UNIQUE NOT NULL,
                             full_name VARCHAR(255) NOT NULL,
                             is_active BOOLEAN DEFAULT TRUE,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 2. DEVICE MANAGEMENT (FR76)
-- ==========================================

-- Tracks hardware assets
CREATE TABLE devices (
                         device_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         serial_number VARCHAR(100) UNIQUE NOT NULL,
                         device_model VARCHAR(100) NOT NULL,
                         firmware_version VARCHAR(50),
                         status VARCHAR(20) DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'ASSIGNED', 'MAINTENANCE', 'RETIRED')),
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tracks assignments for medical auditing (FR77)
CREATE TABLE device_assignments (
                                    assignment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    device_id UUID REFERENCES devices(device_id) ON DELETE CASCADE,
                                    patient_id UUID NOT NULL, -- Reference to the Patient Context
                                    institution_id UUID REFERENCES institutions(institution_id),
                                    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                    unassigned_at TIMESTAMP WITH TIME ZONE,
                                    is_current BOOLEAN DEFAULT TRUE
);

-- ==========================================
-- 3. OPTIMIZATION INDEXES (FR77 / NFR36)
-- ==========================================

-- Tenant isolation index
CREATE INDEX idx_staff_institution ON staff_users(institution_id);

-- Audit log lookup performance
CREATE INDEX idx_assignments_patient ON device_assignments(patient_id);

-- Current assignment filter
CREATE INDEX idx_assignments_active ON device_assignments(device_id) WHERE is_current = TRUE;

-- ==========================================
-- 4. INITIAL SEEDING (Optional for Testing)
-- ==========================================
INSERT INTO roles (role_name, description) VALUES
                                               ('ADMIN', 'System administrator with full access'),
                                               ('DOCTOR', 'Medical professional authorized to start investigations'),
                                               ('NURSE', 'Support staff with read-only access to telemetry');