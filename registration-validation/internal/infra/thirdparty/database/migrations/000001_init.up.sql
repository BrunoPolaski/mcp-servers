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

CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    last_verified_at TIMESTAMP
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

CREATE INDEX idx_compliance_checks_person_id ON compliance_checks(person_id);
CREATE INDEX idx_compliance_checks_check_type ON compliance_checks(check_type);
CREATE INDEX idx_compliance_checks_check_date ON compliance_checks(check_date);
CREATE INDEX idx_compliance_checks_is_pep ON compliance_checks(is_pep);
CREATE INDEX idx_compliance_checks_on_sanctions_list ON compliance_checks(on_sanctions_list);

