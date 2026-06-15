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
    document_validated BOOLEAN DEFAULT FALSE,
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    biometric_validated BOOLEAN DEFAULT FALSE,
    receita_federal_status VARCHAR(50)
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

CREATE TABLE credit_scores (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL,
    score INT NOT NULL,
    score_date TIMESTAMP NOT NULL,
    score_model VARCHAR(50) NOT NULL,
    score_reason TEXT,
    payment_history INT NOT NULL,
    credit_usage INT NOT NULL,
    credit_age INT NOT NULL,
    credit_mix INT NOT NULL,
    recent_inquiries INT NOT NULL,
    risk_level VARCHAR(50) NOT NULL,
    default_probability DOUBLE PRECISION NOT NULL,
    UNIQUE (person_id, score_date)
);

CREATE TABLE financial_profiles (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL,
    profile_date TIMESTAMP NOT NULL,
    declared_monthly_income DOUBLE PRECISION,
    estimated_monthly_income DOUBLE PRECISION,
    income_source VARCHAR(100),
    total_assets DOUBLE PRECISION,
    real_estate_value DOUBLE PRECISION,
    vehicles_value DOUBLE PRECISION,
    investments_value DOUBLE PRECISION,
    total_liabilities DOUBLE PRECISION,
    total_monthly_payments DOUBLE PRECISION,
    debt_to_income_ratio DOUBLE PRECISION,
    available_credit DOUBLE PRECISION,
    credit_utilization DOUBLE PRECISION,
    banking_relationships INT,
    account_age_average INT,
    has_checking_account BOOLEAN DEFAULT FALSE,
    has_savings_account BOOLEAN DEFAULT FALSE,
    has_investment_account BOOLEAN DEFAULT FALSE,
    UNIQUE (person_id, profile_date)
);

CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    credit_score_id BIGINT,
    financial_profile_id BIGINT,
    last_verified_at TIMESTAMP,
    consent_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    consent_granted_at TIMESTAMP
);

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

CREATE TABLE credit_accounts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    account_type VARCHAR(50) NOT NULL,
    creditor VARCHAR(255) NOT NULL,
    creditor_document VARCHAR(14),
    account_number VARCHAR(100),
    opened_date TIMESTAMP NOT NULL,
    closed_date TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    credit_limit DOUBLE PRECISION,
    current_balance DOUBLE PRECISION NOT NULL,
    available_credit DOUBLE PRECISION,
    original_amount DOUBLE PRECISION,
    remaining_amount DOUBLE PRECISION,
    interest_rate DOUBLE PRECISION,
    monthly_payment DOUBLE PRECISION,
    payment_due_day INT,
    number_of_payments INT,
    remaining_payments INT,
    payment_status VARCHAR(50),
    days_late INT DEFAULT 0,
    highest_days_late INT DEFAULT 0,
    times_late_30_days INT DEFAULT 0,
    times_late_60_days INT DEFAULT 0,
    times_late_90_days INT DEFAULT 0,
    last_reported_date TIMESTAMP NOT NULL
);

CREATE TABLE credit_inquiries (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    inquiry_date TIMESTAMP NOT NULL,
    inquiry_type VARCHAR(50) NOT NULL,
    creditor VARCHAR(255) NOT NULL,
    creditor_document VARCHAR(14),
    purpose VARCHAR(100),
    amount DOUBLE PRECISION,
    result VARCHAR(50)
);

CREATE TABLE debts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    debt_type VARCHAR(100) NOT NULL,
    creditor VARCHAR(255) NOT NULL,
    creditor_document VARCHAR(14),
    original_amount DOUBLE PRECISION NOT NULL,
    current_amount DOUBLE PRECISION NOT NULL,
    interest_rate DOUBLE PRECISION,
    fees DOUBLE PRECISION,
    origin_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    in_collection BOOLEAN DEFAULT FALSE,
    collection_date TIMESTAMP,
    collection_agency VARCHAR(255),
    settlement_amount DOUBLE PRECISION,
    settlement_date TIMESTAMP
);

CREATE TABLE payment_histories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    credit_account_id BIGINT REFERENCES credit_accounts(id),
    debt_id BIGINT REFERENCES debts(id),
    payment_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    amount_due DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) NOT NULL,
    days_late INT DEFAULT 0
);

CREATE TABLE negative_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    record_type VARCHAR(100) NOT NULL,
    creditor VARCHAR(255) NOT NULL,
    creditor_document VARCHAR(14),
    amount DOUBLE PRECISION NOT NULL,
    inclusion_date TIMESTAMP NOT NULL,
    contract_number VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    removal_date TIMESTAMP,
    removal_reason VARCHAR(255),
    process_number VARCHAR(100),
    notary VARCHAR(255),
    is_disputed BOOLEAN DEFAULT FALSE,
    dispute_date TIMESTAMP,
    dispute_reason TEXT
);

CREATE TABLE employment_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    employer_name VARCHAR(255) NOT NULL,
    employer_document VARCHAR(14),
    job_title VARCHAR(255),
    employment_type VARCHAR(50),
    salary DOUBLE PRECISION,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    is_current BOOLEAN DEFAULT FALSE,
    verification_status VARCHAR(50) DEFAULT 'unverified',
    data_source VARCHAR(100)
);

CREATE TABLE income_declarations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    declaration_date TIMESTAMP NOT NULL,
    income_type VARCHAR(100) NOT NULL,
    monthly_amount DOUBLE PRECISION NOT NULL,
    yearly_amount DOUBLE PRECISION,
    source VARCHAR(255),
    verified BOOLEAN DEFAULT FALSE,
    verified_by VARCHAR(100),
    proof_file_id BIGINT REFERENCES files(id)
);

CREATE TABLE legal_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    record_type VARCHAR(100) NOT NULL,
    process_number VARCHAR(100),
    court VARCHAR(255),
    filing_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    amount DOUBLE PRECISION,
    description TEXT,
    resolution TEXT,
    resolution_date TIMESTAMP
);

CREATE TABLE compliance_checks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    check_type VARCHAR(100) NOT NULL,
    check_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    details JSONB,
    is_pep BOOLEAN DEFAULT FALSE,
    pep_details TEXT,
    on_sanctions_list BOOLEAN DEFAULT FALSE,
    sanctions_details TEXT,
    valid_until TIMESTAMP
);

CREATE TABLE fraud_alerts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    alert_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    detected_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    resolved_date TIMESTAMP,
    resolved_by VARCHAR(255),
    notes TEXT
);

CREATE TABLE risk_assessments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    assessment_date TIMESTAMP NOT NULL,
    assessment_type VARCHAR(100) NOT NULL,
    risk_score INT NOT NULL,
    risk_level VARCHAR(50) NOT NULL,
    risk_factors JSONB,
    recommendation TEXT,
    model_version VARCHAR(50)
);

CREATE TABLE person_relationships (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    related_person_id BIGINT NOT NULL REFERENCES persons(id),
    relation_type VARCHAR(100) NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    verification_status VARCHAR(50) DEFAULT 'unverified'
);

CREATE TABLE data_sources (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    source_name VARCHAR(255) NOT NULL UNIQUE,
    source_type VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_sync_date TIMESTAMP,
    reliability_score INT
);

CREATE TABLE person_data_sources (
    person_id BIGINT NOT NULL REFERENCES persons(id),
    data_source_id BIGINT NOT NULL REFERENCES data_sources(id),
    PRIMARY KEY (person_id, data_source_id)
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

CREATE INDEX idx_persons_personal_information_id ON persons(personal_information_id);
CREATE INDEX idx_persons_credit_score_id ON persons(credit_score_id);
CREATE INDEX idx_persons_financial_profile_id ON persons(financial_profile_id);

CREATE INDEX idx_credit_scores_score ON credit_scores(score);

CREATE INDEX idx_financial_profiles_declared_monthly_income ON financial_profiles(declared_monthly_income);
CREATE INDEX idx_financial_profiles_estimated_monthly_income ON financial_profiles(estimated_monthly_income);
CREATE INDEX idx_financial_profiles_total_liabilities ON financial_profiles(total_liabilities);
CREATE INDEX idx_financial_profiles_debt_to_income_ratio ON financial_profiles(debt_to_income_ratio);
CREATE INDEX idx_financial_profiles_credit_utilization ON financial_profiles(credit_utilization);

CREATE INDEX idx_credit_accounts_person_id ON credit_accounts(person_id);
CREATE INDEX idx_credit_accounts_account_type ON credit_accounts(account_type);
CREATE INDEX idx_credit_accounts_status ON credit_accounts(status);
CREATE INDEX idx_credit_accounts_credit_limit ON credit_accounts(credit_limit);
CREATE INDEX idx_credit_accounts_current_balance ON credit_accounts(current_balance);
CREATE INDEX idx_credit_accounts_payment_status ON credit_accounts(payment_status);
CREATE INDEX idx_credit_accounts_days_late ON credit_accounts(days_late);

CREATE INDEX idx_credit_inquiries_person_id ON credit_inquiries(person_id);
CREATE INDEX idx_credit_inquiries_inquiry_date ON credit_inquiries(inquiry_date);

CREATE INDEX idx_payment_histories_person_id ON payment_histories(person_id);
CREATE INDEX idx_payment_histories_credit_account_id ON payment_histories(credit_account_id);
CREATE INDEX idx_payment_histories_debt_id ON payment_histories(debt_id);
CREATE INDEX idx_payment_histories_payment_date ON payment_histories(payment_date);
CREATE INDEX idx_payment_histories_status ON payment_histories(status);
CREATE INDEX idx_payment_histories_days_late ON payment_histories(days_late);

CREATE INDEX idx_debts_person_id ON debts(person_id);
CREATE INDEX idx_debts_debt_type ON debts(debt_type);
CREATE INDEX idx_debts_creditor ON debts(creditor);
CREATE INDEX idx_debts_current_amount ON debts(current_amount);
CREATE INDEX idx_debts_due_date ON debts(due_date);
CREATE INDEX idx_debts_status ON debts(status);
CREATE INDEX idx_debts_in_collection ON debts(in_collection);

CREATE INDEX idx_negative_records_person_id ON negative_records(person_id);
CREATE INDEX idx_negative_records_record_type ON negative_records(record_type);
CREATE INDEX idx_negative_records_status ON negative_records(status);

CREATE INDEX idx_employment_records_person_id ON employment_records(person_id);
CREATE INDEX idx_employment_records_is_current ON employment_records(is_current);

CREATE INDEX idx_income_declarations_person_id ON income_declarations(person_id);
CREATE INDEX idx_income_declarations_declaration_date ON income_declarations(declaration_date);

CREATE INDEX idx_legal_records_person_id ON legal_records(person_id);
CREATE INDEX idx_legal_records_record_type ON legal_records(record_type);
CREATE INDEX idx_legal_records_process_number ON legal_records(process_number);

CREATE INDEX idx_compliance_checks_person_id ON compliance_checks(person_id);
CREATE INDEX idx_compliance_checks_check_type ON compliance_checks(check_type);
CREATE INDEX idx_compliance_checks_check_date ON compliance_checks(check_date);
CREATE INDEX idx_compliance_checks_is_pep ON compliance_checks(is_pep);
CREATE INDEX idx_compliance_checks_on_sanctions_list ON compliance_checks(on_sanctions_list);

CREATE INDEX idx_fraud_alerts_person_id ON fraud_alerts(person_id);
CREATE INDEX idx_fraud_alerts_alert_type ON fraud_alerts(alert_type);
CREATE INDEX idx_fraud_alerts_severity ON fraud_alerts(severity);
CREATE INDEX idx_fraud_alerts_detected_date ON fraud_alerts(detected_date);

CREATE INDEX idx_risk_assessments_person_id ON risk_assessments(person_id);
CREATE INDEX idx_risk_assessments_assessment_date ON risk_assessments(assessment_date);
CREATE INDEX idx_risk_assessments_risk_score ON risk_assessments(risk_score);
CREATE INDEX idx_risk_assessments_risk_level ON risk_assessments(risk_level);

CREATE INDEX idx_person_relationships_person_id ON person_relationships(person_id);
CREATE INDEX idx_person_relationships_related_person_id ON person_relationships(related_person_id);

CREATE INDEX idx_data_sources_reliability_score ON data_sources(reliability_score);