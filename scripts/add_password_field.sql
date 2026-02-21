-- Add password field to users table
-- Run this migration after adding JWT auth

USE user_service;

-- Temporarily disable safe update mode
SET SQL_SAFE_UPDATES = 0;

-- Add password column to users table
ALTER TABLE users ADD COLUMN password VARCHAR(255) NOT NULL DEFAULT '' AFTER name;

-- Update existing users with a default hashed password (change this!)
-- Note: This is a bcrypt hash of "changeme123" - users should change their password
UPDATE users SET password = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' WHERE password = '';

-- Re-enable safe update mode
SET SQL_SAFE_UPDATES = 1;
