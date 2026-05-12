CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    zip_code VARCHAR(16) NOT NULL,
    state VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL,
    number VARCHAR(16) NOT NULL,
    complement VARCHAR(255)
);

CREATE TABLE api_keys (
    UUID CHAR(36) PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    original_name VARCHAR(255),
    name VARCHAR(255),
    url VARCHAR(255),
    mime_type VARCHAR(50)
);


CREATE TABLE personal_informations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    document VARCHAR(11) NOT NULL UNIQUE,
    address_id BIGINT REFERENCES addresses(id),
    file_id BIGINT REFERENCES files(id)
);

CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id)
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
    personal_information_id BIGINT NOT NULL REFERENCES personal_informations(id),
    business_id BIGINT NOT NULL REFERENCES businesses(id)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    email VARCHAR(255) NOT NULL UNIQUE,
    user_type VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    owner_id BIGINT REFERENCES owners(id),
    person_id BIGINT REFERENCES persons(id),
    analyst_id BIGINT REFERENCES analysts(id),
    admin_id BIGINT REFERENCES admins(id)
);

CREATE TABLE tokens (
    UUID CHAR(36) PRIMARY KEY,
    api_key CHAR(36),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);