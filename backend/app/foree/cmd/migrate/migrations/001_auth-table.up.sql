CREATE TABLE IF NOT EXISTS users(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32),
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

CREATE TABLE IF NOT EXISTS user_group(
    `id` SERIAL PRIMARY KEY,
    `role_group`  VARCHAR(128),
    `transaction_limit_group` VARCHAR(128),
    `owner_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS email_passwd(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32),
    `email` VARCHAR(255) UNIQUE KEY NOT NULL,
)