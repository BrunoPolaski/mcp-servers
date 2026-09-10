CREATE TABLE document_validations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    validation_date TIMESTAMP NOT NULL,
    document_number VARCHAR(14) NOT NULL,
    document_type VARCHAR(10) NOT NULL,
    receita_federal_status VARCHAR(50) NOT NULL,
    is_valid BOOLEAN DEFAULT FALSE,
    name_matches BOOLEAN DEFAULT FALSE,
    birth_date_matches BOOLEAN DEFAULT FALSE,
    biometric_validated BOOLEAN DEFAULT FALSE,
    source VARCHAR(100) DEFAULT 'receita_federal',
    raw_response JSONB
);

CREATE TABLE fiscal_regularities (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    check_date TIMESTAMP NOT NULL,
    has_debts BOOLEAN DEFAULT FALSE,
    cnd_status VARCHAR(50) NOT NULL,
    cnd_number VARCHAR(100),
    cnd_issue_date TIMESTAMP,
    cnd_valid_until TIMESTAMP,
    pending_issues JSONB
);

CREATE TABLE employment_link_validations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    person_id BIGINT NOT NULL REFERENCES persons(id),
    validation_date TIMESTAMP NOT NULL,
    employer_name VARCHAR(255) NOT NULL,
    employer_document VARCHAR(14),
    employment_type VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    source VARCHAR(100) DEFAULT 'eSocial',
    verified BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_document_validations_person_id ON document_validations(person_id);
CREATE INDEX idx_document_validations_validation_date ON document_validations(validation_date);
CREATE INDEX idx_document_validations_document_number ON document_validations(document_number);
CREATE INDEX idx_document_validations_is_valid ON document_validations(is_valid);

CREATE INDEX idx_fiscal_regularities_person_id ON fiscal_regularities(person_id);
CREATE INDEX idx_fiscal_regularities_check_date ON fiscal_regularities(check_date);
CREATE INDEX idx_fiscal_regularities_has_debts ON fiscal_regularities(has_debts);

CREATE INDEX idx_employment_link_validations_person_id ON employment_link_validations(person_id);
CREATE INDEX idx_employment_link_validations_validation_date ON employment_link_validations(validation_date);
CREATE INDEX idx_employment_link_validations_status ON employment_link_validations(status);
CREATE INDEX idx_employment_link_validations_verified ON employment_link_validations(verified);
