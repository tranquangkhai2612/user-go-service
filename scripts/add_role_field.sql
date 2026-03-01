-- Add role field to users table
-- Run this migration to add role-based access control

USE user_service;

-- Temporarily disable safe update mode
SET SQL_SAFE_UPDATES = 0;

-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER name;

-- Create an admin user (optional - update email/password as needed)
-- Password is bcrypt hash of "admin123"
INSERT INTO users (id, email, name, role, password, created_at, updated_at) 
VALUES (
    UUID(), 
    'admin@example.com', 
    'System Admin', 
    'admin',
    '$2a$10$qQ8xqF7K7F5X5X5X5X5X5uN9YxZ8X5X5X5X5X5X5X5X5X5X5X5X5X',
    NOW(), 
    NOW()
) ON DUPLICATE KEY UPDATE role = 'admin';

-- Re-enable safe update mode
SET SQL_SAFE_UPDATES = 1;
