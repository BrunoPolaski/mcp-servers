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

CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    customer_relationship_id BIGINT,
    last_verified_at TIMESTAMP
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

CREATE TABLE customer_relationships (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL, -- sem FK: carregada antes de persons (persons.customer_relationship_id aponta para cá), espelhando bank_account_profiles do open-finance
    customer_since TIMESTAMP NOT NULL,
    relationship_months INTEGER NOT NULL DEFAULT 0,
    segment VARCHAR(50),
    branch VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    churn_risk VARCHAR(50),
    internal_score INTEGER
);
CREATE UNIQUE INDEX idx_customer_relationships_person_id ON customer_relationships(person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_relationships_months ON customer_relationships(relationship_months);
CREATE INDEX idx_customer_relationships_active ON customer_relationships(is_active);
CREATE INDEX idx_customer_relationships_score ON customer_relationships(internal_score);

CREATE TABLE contracted_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    product_type VARCHAR(50) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    contract_number VARCHAR(100),
    contracted_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    balance DOUBLE PRECISION,
    monthly_value DOUBLE PRECISION
);
CREATE INDEX idx_contracted_products_person_id ON contracted_products(person_id);
CREATE INDEX idx_contracted_products_type ON contracted_products(product_type);
CREATE INDEX idx_contracted_products_number ON contracted_products(contract_number);
CREATE INDEX idx_contracted_products_status ON contracted_products(status);

CREATE TABLE internal_payment_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    contracted_product_id BIGINT REFERENCES contracted_products(id),
    reference_month TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    payment_date TIMESTAMP,
    amount_due DOUBLE PRECISION NOT NULL,
    amount_paid DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) NOT NULL,
    days_late INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_internal_payment_records_person_id ON internal_payment_records(person_id);
CREATE INDEX idx_internal_payment_records_product_id ON internal_payment_records(contracted_product_id);
CREATE INDEX idx_internal_payment_records_ref_month ON internal_payment_records(reference_month);
CREATE INDEX idx_internal_payment_records_status ON internal_payment_records(status);
CREATE INDEX idx_internal_payment_records_days_late ON internal_payment_records(days_late);

CREATE TABLE pre_approved_limits (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    product_type VARCHAR(50) NOT NULL,
    approved_amount DOUBLE PRECISION NOT NULL,
    interest_rate DOUBLE PRECISION,
    calculated_date TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    policy_version VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX idx_pre_approved_limits_person_id ON pre_approved_limits(person_id);
CREATE INDEX idx_pre_approved_limits_type ON pre_approved_limits(product_type);
CREATE INDEX idx_pre_approved_limits_valid_until ON pre_approved_limits(valid_until);
CREATE INDEX idx_pre_approved_limits_active ON pre_approved_limits(is_active);

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
CREATE INDEX idx_persons_customer_relationship_id ON persons(customer_relationship_id);

CREATE INDEX idx_income_declarations_person_id ON income_declarations(person_id);
CREATE INDEX idx_income_declarations_declaration_date ON income_declarations(declaration_date);
