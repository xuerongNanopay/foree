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

CREATE TABLE IF NOT EXISTS user_extra(
    `id` SERIAL PRIMARY KEY,
    `pob` VARCHAR(64),
    `cor` VARCHAR(64),
    `nationality` VARCHAR(64),
    `occupation_category` VARCHAR(64),
    `occupation_name` VARCHAR(128),
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_identifications(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32),
    `type` VARCHAR(32),
    `value` VARCHAR(64),
    `link1` VARCHAR(256),
    `link2` VARCHAR(256),
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS referral(
    `id` SERIAL PRIMARY KEY,
    `referral_type` VARCHAR(32),
    `referral_value` VARCHAR(256),
    `referral_code` VARCHAR(256),
    `referrer_id` BIGINT UNSIGNED NOT NULL,
    `referee_id` BIGINT UNSIGNED,
    `accept_at` DATETIME,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (referrer_id) REFERENCES users(id),
    FOREIGN KEY (referee_id) REFERENCES users(id)
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

CREATE TABLE IF NOT EXISTS foree_tx(
    `id` SERIAL PRIMARY KEY,
    `type` VARCHAR(64),
    `status` VARCHAR(32),
    `cin_acc_id` BIGINT UNSIGNED NOT NULL,
    `cout_acc_id` BIGINT UNSIGNED NOT NULL,
    `rate` DECIMAL(7, 2) NOT NULL,
    `src_amount` DECIMAL(8, 2) NOT NULL,
    `src_currency` CHAR(3) NOT NULL,
    `dest_amount` DECIMAL(8, 2) NOT NULL,
    `dest_currency` CHAR(3) NOT NULL,
    `total_fee_amount` DECIMAL(8, 2) NOT NULL,
    `total_fee_currency` CHAR(3) NOT NULL,
    `total_reward_amount` DECIMAL(8, 2) NOT NULL,
    `total_reward_currency` CHAR(3) NOT NULL,
    `total_amount` DECIMAL(8, 2) NOT NULL,
    `total_currency` CHAR(3) NOT NULL,
    `cur_stage` VARCHAR(64) NOT NULL,
    `cur_stage_status` VARCHAR(32) NOT NULL,
    `transaction_purpose` VARCHAR(256) NOT NULL,
    `conclusion` VARCHAR(256) NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS idm_tx(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(64),
    `ip` VARCHAR(16),
    `user_agent` VARCHAR(256),
    `idm_reference` VARCHAR(64),
    `idm_result` VARCHAR(64),
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS idm_compliance(
    `id` SERIAL PRIMARY KEY,
    `idm_tx_id` BIGINT UNSIGNED NOT NULL,
    `idm_http_status_code` int,
    `idm_result` VARCHAR(64),
    `request_json` VARCHAR(1024),
    `response_json` VARCHAR(4096),
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);
