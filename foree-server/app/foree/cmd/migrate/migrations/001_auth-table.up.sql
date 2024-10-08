CREATE TABLE IF NOT EXISTS users(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `first_name` VARCHAR(256) DEFAULT '',
    `middle_name` VARCHAR(256) DEFAULT '',
    `last_name` VARCHAR(256) DEFAULT '',
    `age` TINYINT UNSIGNED,
    `dob` DATE,
    `address1` VARCHAR(256) DEFAULT '',
    `address2` VARCHAR(256) DEFAULT '',
    `city` VARCHAR(64) DEFAULT '',
    `province` VARCHAR(5) DEFAULT '',
    `country` VARCHAR(2) DEFAULT '',
    `postal_code` VARCHAR(16) DEFAULT '',
    `phone_number` VARCHAR(32) DEFAULT '',
    `email` VARCHAR(256) UNIQUE KEY,
    `avatar_url` VARCHAR(256) DEFAULT '',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_setting(
    `id` SERIAL PRIMARY KEY,
    `is_in_app_notification_enable` BOOL DEFAULT TRUE,
    `is_push_notification_enable` BOOL DEFAULT TRUE,
    `is_email_notifications_enable` BOOL DEFAULT TRUE,
    `owner_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
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
    `pob` VARCHAR(2) DEFAULT '',
    `cor` VARCHAR(2) DEFAULT '',
    `nationality` VARCHAR(2) DEFAULT '',
    `occupation_category` VARCHAR(64) DEFAULT '',
    `occupation_name` VARCHAR(128) DEFAULT '',
    `owner_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_identifications(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) DEFAULT '',
    `type` VARCHAR(32) DEFAULT '',
    `value` VARCHAR(64) DEFAULT '',
    `link1` VARCHAR(256) DEFAULT '',
    `link2` VARCHAR(256) DEFAULT '',
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS referral(
    `id` SERIAL PRIMARY KEY,
    `referral_type` VARCHAR(32) DEFAULT '',
    `referral_value` VARCHAR(256) DEFAULT '',
    `referral_code` VARCHAR(256) DEFAULT '',
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
    `username` VARCHAR(256) UNIQUE KEY,
    `passwd` VARCHAR(32) NOT NULL,
    `verify_code` VARCHAR(32) DEFAULT '',
    `verify_code_expired_at` DATETIME,
    `login_attempts` INT DEFAULT 0,
    `retrieve_token` VARCHAR(128) DEFAULT '',
    `retrieve_token_expired_at` DATETIME,
    `latest_forget_passwd_updated_at` DATETIME,
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
    `middle_name` VARCHAR(256) DEFAULT '',
    `last_name` VARCHAR(256) NOT NULL,
    `address1` VARCHAR(256) NOT NULL,
    `address2` VARCHAR(256) DEFAULT '',
    `city` VARCHAR(64) NOT NULL,
    `province` CHAR(5) NOT NULL,
    `country` CHAR(2) NOT NULL,
    `postal_code` VARCHAR(16) DEFAULT '',
    `phone_number` VARCHAR(32) DEFAULT '',
    `institution_name` VARCHAR(128) DEFAULT '',
    `branch_number` VARCHAR(64) DEFAULT '',
    `account_number`VARCHAR(128) DEFAULT '',
    `account_hash`VARCHAR(256) DEFAULT '',
    `relationship_to_contact`VARCHAR(128) DEFAULT '',
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
    `middle_name` VARCHAR(256) DEFAULT '',
    `last_name` VARCHAR(256) NOT NULL,
    `email` VARCHAR(256) NOT NULL,
    `latest_activity_at` DATETIME,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS tx_limit(
    `name` VARCHAR(128) NOT NULL UNIQUE,
    `limit_group` VARCHAR(64) NOT NULL UNIQUE,
    `min_amount` DECIMAL(10, 2) NOT NULL,
    `min_currency` CHAR(3) NOT NULL,
    `max_amount` DECIMAL(10, 2) NOT NULL,
    `max_currency` CHAR(3) NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
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
    `description` VARCHAR(256) DEFAULT '',
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
    `description` VARCHAR(256) DEFAULT '',
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
    `description` VARCHAR(256) DEFAULT '',
    `amount` DECIMAL(7, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `applied_transaction_id` BIGINT UNSIGNED DEFAULT 0,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `expire_at` DATETIME,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_extra(
    `id` SERIAL PRIMARY KEY,
    `pob` VARCHAR(64) DEFAULT '',
    `cor` VARCHAR(64) DEFAULT '',
    `nationality` VARCHAR(64) DEFAULT '',
    `occupation_category` VARCHAR(64) DEFAULT '',
    `occupation_name` VARCHAR(128) DEFAULT '',
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_identifications(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(32) NOT NULL,
    `type` VARCHAR(32) DEFAULT '',
    `value` VARCHAR(64) DEFAULT '',
    `link1` VARCHAR(256) DEFAULT '',
    `link2` VARCHAR(256) DEFAULT '',
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS referral(
    `id` SERIAL PRIMARY KEY,
    `referral_type` VARCHAR(32) DEFAULT '',
    `referral_value` VARCHAR(256) DEFAULT '',
    `referral_code` VARCHAR(256) DEFAULT '',
    `referrer_id` BIGINT UNSIGNED NOT NULL,
    `referee_id` BIGINT UNSIGNED,
    `accept_at` DATETIME,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (referrer_id) REFERENCES users(id),
    FOREIGN KEY (referee_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS foree_tx(
    `id` SERIAL PRIMARY KEY,
    `type` VARCHAR(64) NOT NULL,
    `cin_acc_id` BIGINT UNSIGNED NOT NULL,
    `cout_acc_id` BIGINT UNSIGNED NOT NULL,
    `limit_reference` VARCHAR(64) NOT NULL,
    `rate` DECIMAL(7, 2) NOT NULL,
    `src_amount` DECIMAL(11, 2) NOT NULL,
    `src_currency` CHAR(3) NOT NULL,
    `dest_amount` DECIMAL(11, 2) NOT NULL,
    `dest_currency` CHAR(3) NOT NULL,
    `total_fee_amount` DECIMAL(11, 2) NOT NULL,
    `total_fee_currency` CHAR(3) NOT NULL,
    `total_reward_amount` DECIMAL(11, 2) NOT NULL,
    `total_reward_currency` CHAR(3) NOT NULL,
    `total_amount` DECIMAL(11, 2) NOT NULL,
    `total_currency` CHAR(3) NOT NULL,
    `stage` VARCHAR(64) NOT NULL,
    `transaction_purpose` VARCHAR(256) NOT NULL,
    `conclusion` VARCHAR(256) NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS idm_tx(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(64) NOT NULL,
    `ip` VARCHAR(16) DEFAULT '',
    `user_agent` VARCHAR(256) DEFAULT '',
    `idm_reference` VARCHAR(64) DEFAULT '',
    `idm_result` VARCHAR(64) DEFAULT '',
    `parent_tx_id` BIGINT UNSIGNED NOT NULL UNIQUE,
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
    `idm_result` VARCHAR(64) DEFAULT '',
    `request_json` VARCHAR(1024),
    `response_json` VARCHAR(4096),
    `parent_tx_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS interac_ci_tx(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(64) NOT NULL,
    `cash_in_acc_id` BIGINT UNSIGNED NOT NULL,
    `amount` DECIMAL(11, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `scotia_payment_id` VARCHAR(128) DEFAULT '',
    `scotia_status` VARCHAR(64) DEFAULT '',
    `scotia_clearing_reference` VARCHAR(128) DEFAULT '',
    `payment_url` VARCHAR(256) DEFAULT '',
    `end_to_end_id` VARCHAR(128) DEFAULT '',
    `parent_tx_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS foree_refund_tx(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(64) NOT NULL,
    `refund_interac_acc_id` BIGINT UNSIGNED NOT NULL,
    `refund_amount` DECIMAL(11, 2) NOT NULL,
    `refund_currency` CHAR(3) NOT NULL,
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS nbp_co_tx(
    `id` SERIAL PRIMARY KEY,
    `status` VARCHAR(64) NOT NULL,
    `mode` VARCHAR(64) DEFAULT '',
    `amount` DECIMAL(11, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `nbp_reference` VARCHAR(128) DEFAULT '',
    `cash_out_acc_id` BIGINT UNSIGNED NOT NULL,
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS tx_summary(
    `id` SERIAL PRIMARY KEY,
    `summary` VARCHAR(256) NOT NULL,
    `type` VARCHAR(64) NOT NULL,
    `status` VARCHAR(64) NOT NULL,
    `rate` VARCHAR(64) NOT NULL,
    `payment_url` VARCHAR(256) DEFAULT '',
    `src_acc_id` BIGINT UNSIGNED NOT NULL,
    `dest_acc_id` BIGINT UNSIGNED NOT NULL,
    `src_acc_summary` VARCHAR(128) NOT NULL,
    `src_amount` DECIMAL(11, 2) NOT NULL,
    `src_currency` CHAR(3) NOT NULL,
    `dest_acc_summary` VARCHAR(128) NOT NULL,
    `total_amount` DECIMAL(11, 2) NOT NULL,
    `total_currency` CHAR(3) NOT NULL,
    `fee_amount` DECIMAL(11, 2) NOT NULL,
    `fee_currency` CHAR(3) NOT NULL,
    `reward_amount` DECIMAL(11, 2) NOT NULL,
    `dest_currency` CHAR(3) NOT NULL,
    `dest_amount` DECIMAL(11, 2) NOT NULL,
    `reward_currency` CHAR(3) NOT NULL,
    `nbp_reference` VARCHAR(128) DEFAULT '',
    `is_cancel_allowed` BOOLEAN DEFAULT FALSE,
    `parent_tx_id` BIGINT UNSIGNED NOT NULL,
    `owner_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (parent_tx_id) REFERENCES foree_tx(id)
);

CREATE TABLE IF NOT EXISTS promotion(
    `name` VARCHAR(128) NOT NULL PRIMARY KEY,
    `description` VARCHAR(256) DEFAULT '',
    `amount` DECIMAL(11, 2) NOT NULL,
    `currency` CHAR(3) NOT NULL,
    `is_enable` BOOL NOT NULL,
    `start_time` DATETIME,
    `end_time` DATETIME,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_tab_email_passwd_col_email ON email_passwd(email);
CREATE INDEX idx_tab_user_extra_col_owner_id ON user_extra(owner_id);
CREATE INDEX idx_tab_user_setting_col_owner_id ON user_setting(owner_id);
CREATE INDEX idx_tab_user_group_col_owner_id ON user_group(owner_id);
CREATE INDEX idx_tab_user_identifications_col_owner_id ON user_identifications(owner_id);
CREATE INDEX idx_tab_contact_accounts_col_latest_activity_at ON contact_accounts(latest_activity_at);
CREATE INDEX idx_tab_contact_accounts_col_latest_owner_id ON contact_accounts(owner_id);
CREATE INDEX idx_tab_interac_accounts_col_latest_owner_id ON interac_accounts(owner_id);
CREATE INDEX idx_tab_daily_tx_limit_col_reference ON daily_tx_limit(reference);
CREATE INDEX idx_tab_fees_col_name ON fees(name);
CREATE INDEX idx_tab_fees_col_fee_group ON fees(fee_group);
CREATE INDEX idx_tab_fee_joint_col_parent_tx_id ON fee_joint(parent_tx_id);
CREATE INDEX idx_tab_promotion_col_name ON promotion(name);
CREATE INDEX idx_tab_referral_col_referrer_id ON referral(referrer_id);
CREATE INDEX idx_tab_referral_col_referee_id ON referral(referee_id);
CREATE INDEX idx_tab_referral_col_referral_code ON referral(referral_code);
CREATE INDEX idx_tab_rewards_col_owner_id ON rewards(owner_id);
CREATE INDEX idx_promotion_col_name ON promotion(name);
CREATE INDEX idx_interac_ci_tx_col_parent_tx_id ON interac_ci_tx(parent_tx_id);
CREATE INDEX idx_interac_ci_tx_col_scotia_payment_id ON interac_ci_tx(scotia_payment_id);
CREATE INDEX idx_idm_tx_col_parent_tx_id ON idm_tx(parent_tx_id);
CREATE INDEX idx_tab_nbp_co_tx_col_parent_tx_id ON nbp_co_tx(parent_tx_id);
CREATE INDEX idx_tab_nbp_co_tx_col_nbp_reference ON nbp_co_tx(nbp_reference);
CREATE INDEX idx_tab_foree_refund_tx_col_parent_tx_id ON foree_refund_tx(parent_tx_id);
CREATE INDEX idx_tab_tx_summary_col_parent_tx_id ON tx_summary(parent_tx_id);
CREATE INDEX idx_tab_tx_summary_col_owner_id ON tx_summary(owner_id);
CREATE INDEX idx_tab_tx_summary_col_created_at ON tx_summary(created_at);
CREATE INDEX idx_tab_rewards_col_applied_transaction_id ON rewards(applied_transaction_id);