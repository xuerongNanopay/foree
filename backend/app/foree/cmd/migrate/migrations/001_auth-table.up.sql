CREATE TABLE IF NOT EXISTS users(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `first_name` VARCHAR(256),
    `middle_name` VARCHAR(256),
    `last_name` VARCHAR(256),
    `age` TINYINT UNSIGNED,
    `dob` DATE,
    `address1` VARCHAR(256),
    `address2` VARCHAR(256),
    `city` VARCHAR(64),
    `province` VARCHAR(64),
    `country` VARCHAR(64),
    `postal_code` VARCHAR(16),
    `phone_number` VARCHAR(32),
    `email` VARCHAR(256) UNIQUE KEY,
    `avatar_url` VARCHAR(256),
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_group(
    `id` SERIAL PRIMARY KEY,
    `role_group`  VARCHAR(128) NOT NULL,
    `transaction_limit_group` VARCHAR(128) NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS email_passwd(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `email` VARCHAR(256) UNIQUE KEY NOT NULL,
    `passwd` VARCHAR(32) NOT NULL,
    `verify_code` VARCHAR(32),
    `verify_code_expired_at` DATETIME,
    `retrieve_token` VARCHAR(128),
    `retrieve_token_expired_at` DATETIME,
    `owner_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS role_permission(
    `role_name` VARCHAR(128) NOT NULL,
    `permission` VARCHAR(128) NOT NULL,
    `is_enable` BOOL NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contact_accounts(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `type` VARCHAR(32) NOT NULL,
    `first_name` VARCHAR(256) NOT NULL,
    `middle_name` VARCHAR(256),
    `last_name` VARCHAR(256) NOT NULL,
    `address1` VARCHAR(256) NOT NULL,
    `address2` VARCHAR(256),
    `city` VARCHAR(64) NOT NULL,
    `province` VARCHAR(64) NOT NULL,
    `country` VARCHAR(64) NOT NULL,
    `postal_code` VARCHAR(16),
    `phone_number` VARCHAR(32),
    `institution_name` VARCHAR(128),
    `branch_number` VARCHAR(64),
    `account_number`VARCHAR(128),
    `account_hash`VARCHAR(256),
    `relationship_to_contact`VARCHAR(128),
    `latest_activity_at` DATETIME,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS interac_accounts(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `first_name` VARCHAR(256) NOT NULL,
    `middle_name` VARCHAR(256),
    `last_name` VARCHAR(256) NOT NULL,
    `address1` VARCHAR(256) NOT NULL,
    `address2` VARCHAR(256),
    `city` VARCHAR(64) NOT NULL,
    `province` VARCHAR(64) NOT NULL,
    `country` VARCHAR(64) NOT NULL,
    `postal_code` VARCHAR(16) NOT NULL,
    `phone_number` VARCHAR(32),
    `email` VARCHAR(256) NOT NULL,
    `institution_name` VARCHAR(128),
    `branch_number` VARCHAR(64),
    `account_number`VARCHAR(128),
    `latest_activity_at` DATETIME,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);