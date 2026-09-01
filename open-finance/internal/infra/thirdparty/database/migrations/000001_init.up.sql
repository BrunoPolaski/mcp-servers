CREATE TABLE api_keys (
    uuid CHAR(36) PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    original_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    mime_type VARCHAR(50) NOT NULL
);

CREATE TABLE personal_informations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    full_name VARCHAR(255) NOT NULL,
    mother_name VARCHAR(255),
    birth_date TIMESTAMP NOT NULL,
    gender VARCHAR(20),
    nationality VARCHAR(100) NOT NULL DEFAULT 'Brazilian',
    marital_status VARCHAR(50),
    document VARCHAR(11) NOT NULL UNIQUE,
    rg VARCHAR(20),
    rg_issuer VARCHAR(50),
    rg_issue_date TIMESTAMP,
    voter_id VARCHAR(20),
    work_card VARCHAR(20),
    primary_phone VARCHAR(15) NOT NULL,
    secondary_phone VARCHAR(15),
    email VARCHAR(255),
    alternative_email VARCHAR(255),
    profile_photo_id BIGINT REFERENCES files(id),
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    zip_code VARCHAR(16) NOT NULL,
    state VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    neighborhood VARCHAR(100),
    street VARCHAR(255) NOT NULL,
    number VARCHAR(16) NOT NULL,
    complement VARCHAR(255),
    reference_point VARCHAR(255),
    address_type VARCHAR(50),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    validated_by_post BOOLEAN DEFAULT FALSE,
    risk_score INT,
    is_current BOOLEAN DEFAULT TRUE,
    is_correspondence BOOLEAN DEFAULT FALSE,
    moved_in_date TIMESTAMP NOT NULL,
    moved_out_date TIMESTAMP,
    verification_status VARCHAR(50) DEFAULT 'unverified'
);

CREATE TABLE person_addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    address_id BIGINT NOT NULL REFERENCES addresses(id)
);

CREATE TABLE person_documents (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    file_id BIGINT NOT NULL REFERENCES files(id),
    document_type VARCHAR(100) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    verified_by VARCHAR(255),
    expiration_date TIMESTAMP
);

CREATE TABLE bank_account_profiles (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL,
    profile_date TIMESTAMP NOT NULL,
    banking_relationships INTEGER NOT NULL DEFAULT 0,
    account_age_average INTEGER,
    has_checking_account BOOLEAN NOT NULL DEFAULT FALSE,
    has_savings_account BOOLEAN NOT NULL DEFAULT FALSE,
    has_investment_account BOOLEAN NOT NULL DEFAULT FALSE,
    investments_value NUMERIC(15, 2)
);

CREATE UNIQUE INDEX idx_person_bank_profile_date
    ON bank_account_profiles (person_id, profile_date);

CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    bank_account_profile_id BIGINT,
    last_verified_at TIMESTAMP
);

CREATE INDEX idx_persons_personal_information_id ON persons (personal_information_id);
CREATE INDEX idx_persons_bank_account_profile_id ON persons (bank_account_profile_id);

CREATE TABLE bank_statements (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    institution VARCHAR(255) NOT NULL,
    institution_document VARCHAR(14),
    account_type VARCHAR(50) NOT NULL,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    opening_balance NUMERIC(15, 2) NOT NULL,
    closing_balance NUMERIC(15, 2) NOT NULL,
    total_credits NUMERIC(15, 2) NOT NULL,
    total_debits NUMERIC(15, 2) NOT NULL,
    transaction_count INTEGER DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'BRL'
);

CREATE INDEX idx_bank_statements_person_id ON bank_statements (person_id);
CREATE INDEX idx_bank_statements_period_start ON bank_statements (period_start);
CREATE INDEX idx_bank_statements_period_end ON bank_statements (period_end);

CREATE TABLE cash_flow_analyses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    analysis_date TIMESTAMP NOT NULL,
    period_days INTEGER NOT NULL,
    average_monthly_inflow NUMERIC(15, 2) NOT NULL,
    average_monthly_outflow NUMERIC(15, 2) NOT NULL,
    net_cash_flow NUMERIC(15, 2) NOT NULL,
    inflow_volatility NUMERIC(6, 4),
    negative_balance_days INTEGER DEFAULT 0,
    has_recurring_income BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_cash_flow_analyses_person_id ON cash_flow_analyses (person_id);
CREATE INDEX idx_cash_flow_analyses_analysis_date ON cash_flow_analyses (analysis_date);
CREATE INDEX idx_cash_flow_analyses_net_cash_flow ON cash_flow_analyses (net_cash_flow);
CREATE INDEX idx_cash_flow_analyses_has_recurring_income ON cash_flow_analyses (has_recurring_income);

CREATE TABLE recurring_transactions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    transaction_type VARCHAR(20) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    amount NUMERIC(15, 2) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    counterparty VARCHAR(255),
    first_detected_date TIMESTAMP NOT NULL,
    last_occurrence_date TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_recurring_transactions_person_id ON recurring_transactions (person_id);
CREATE INDEX idx_recurring_transactions_type ON recurring_transactions (transaction_type);
CREATE INDEX idx_recurring_transactions_last_occurrence ON recurring_transactions (last_occurrence_date);
CREATE INDEX idx_recurring_transactions_is_active ON recurring_transactions (is_active);

CREATE TABLE data_sharing_consents (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    consent_id VARCHAR(100) NOT NULL UNIQUE,
    institution VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    scope JSON,
    granted_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP
);

CREATE INDEX idx_data_sharing_consents_person_id ON data_sharing_consents (person_id);
CREATE INDEX idx_data_sharing_consents_status ON data_sharing_consents (status);

CREATE TABLE admins (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id)
);

CREATE TABLE analysts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    email VARCHAR(255) NOT NULL UNIQUE,
    user_type VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    person_id BIGINT REFERENCES persons(id),
    analyst_id BIGINT REFERENCES analysts(id),
    admin_id BIGINT REFERENCES admins(id)
);

CREATE TABLE sessions (
    uuid CHAR(36) PRIMARY KEY,
    api_key CHAR(36),
    user_id BIGINT NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP,
    last_activity TIMESTAMP
);

CREATE INDEX idx_personal_informations_full_name ON personal_informations(full_name);
CREATE INDEX idx_personal_informations_mother_name ON personal_informations(mother_name);
CREATE INDEX idx_personal_informations_birth_date ON personal_informations(birth_date);
CREATE INDEX idx_personal_informations_rg ON personal_informations(rg);
CREATE INDEX idx_personal_informations_voter_id ON personal_informations(voter_id);
CREATE INDEX idx_personal_informations_primary_phone ON personal_informations(primary_phone);
CREATE INDEX idx_personal_informations_email ON personal_informations(email);

CREATE INDEX idx_addresses_zip_code ON addresses(zip_code);
CREATE INDEX idx_addresses_state ON addresses(state);
CREATE INDEX idx_addresses_city ON addresses(city);
CREATE INDEX idx_addresses_risk_score ON addresses(risk_score);
CREATE INDEX idx_addresses_is_current ON addresses(is_current);
