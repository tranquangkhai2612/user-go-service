-- User Service MySQL Database Initialization Script
-- Created: 2026-02-15

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS user_service DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE user_service;

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample data (optional - remove if not needed)
INSERT INTO users (id, email, name, created_at, updated_at) VALUES
    (UUID(), 'john.doe@example.com', 'John Doe', NOW(), NOW()),
    (UUID(), 'jane.smith@example.com', 'Jane Smith', NOW(), NOW()),
    (UUID(), 'bob.wilson@example.com', 'Bob Wilson', NOW(), NOW());

-- Create a read-only user for reporting purposes (optional)
-- CREATE USER IF NOT EXISTS 'user_service_app'@'%' IDENTIFIED BY 'your_secure_password';
-- GRANT SELECT, INSERT, UPDATE, DELETE ON user_service.* TO 'user_service_app'@'%';
-- FLUSH PRIVILEGES;
