CREATE TABLE IF NOT EXISTS users(
    `id` SERIAL PRIMARY KEY,
    `first_name` VARCHAR(255),
    `middle_name` VARCHAR(255),
    `last_name` VARCHAR(255),
    `age` TINYINT UNSIGNED,
    `dob` DATE,
    `address1` VARCHAR(255),
    `address2` VARCHAR(255),
    `city` VARCHAR(64),
    `province` VARCHAR(64),
    `country` VARCHAR(64),
    `postal_code` VARCHAR(16),
    `phone_number` VARCHAR(32),
    `email` VARCHAR(255) UNIQUE KEY,
    `avatar_url` VARCHAR(255),
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);