-- Create the database
CREATE DATABASE IF NOT EXISTS smart_contract_service;

-- Switch to the database
USE smart_contract_service;

-- Create the user_wallet_log table
CREATE TABLE IF NOT EXISTS user_wallet_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nric VARCHAR(15) NOT NULL UNIQUE,
    wallet_address TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    INDEX idx_nric (nric)
    INDEX idx_wallet_address (wallet_address)
);
