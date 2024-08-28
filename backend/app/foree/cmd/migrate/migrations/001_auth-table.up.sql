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
    `fee_group` VARCHAR(128) NOT NULL,
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

CREATE TABLE IF NOT EXISTS daily_tx_limit(
    `id` SERIAL PRIMARY KEY,
    `reference` VARCHAR(64) NOT NULL,
    `used_amount` DECIMAL(10, 2) NOT NULL,
    `used_currency` CHAR(3) NOT NULL,
    `max_amount` DECIMAL(10, 2) NOT NULL,
    `max_currency` CHAR(3) NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS fees(
    `name` VARCHAR(128) NOT NULL UNIQUE PRIMARY KEY,
    `description` VARCHAR(256),
    `fee_group` VARCHAR(128) NOT NULL,
    `type` VARCHAR(64) NOT NULL,
    `condition` VARCHAR(16) NOT NULL,
    `condition_amount` DECIMAL(5, 2) NOT NULL,
    `condition_currency` CHAR(3) NOT NULL,
    `ratio` DECIMAL(5, 2) NOT NULL,
    `is_apply_in_condition_amount_only` BOOL DEFAULT FALSE,
    `is_enable` BOOL DEFAULT TRUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fee_joint(
    `id` SERIAL PRIMARY KEY,
    `fee_name` VARCHAR(128) NOT NULL,
    `description` VARCHAR(256),
    `amount` DECIMAL(7, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    #FOREIGN KEY ON transaction
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS rate(
    `id` CHAR(7) PRIMARY KEY UNIQUE NOT NULL,
    `src_amount` DECIMAL(7, 2) NOT NULL,
    `src_currency` CHAR(3) NOT NULL,
    `dest_amount` DECIMAL(7, 2) NOT NULL,
    `dest_currency` CHAR(3) NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rewards(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `type` VARCHAR(32) NOT NULL,
    `description` VARCHAR(256),
    `amount` DECIMAL(7, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `applied_transaction_id` BIGINT UNSIGNED,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `expire_at` DATETIME,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);