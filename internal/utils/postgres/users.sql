CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT pk_users_id PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL
        CONSTRAINT ck_users_first_name_length CHECK (LENGTH(TRIM(first_name)) > 0),
    middle_name VARCHAR(50)
        CONSTRAINT ck_users_middle_name_length CHECK (LENGTH(TRIM(middle_name)) > 0),
    last_name VARCHAR(50) NOT NULL
        CONSTRAINT ck_users_last_name_length CHECK (LENGTH(TRIM(last_name)) > 0),
    email VARCHAR(50) NOT NULL
        CONSTRAINT uq_users_email UNIQUE,
        CONSTRAINT ck_users_email_length CHECK (LENGTH(TRIM(email)) >= 5),
    address VARCHAR(100) NOT NULL
        CONSTRAINT ck_users_address_length CHECK (LENGTH(TRIM(address)) > 3),
    gender VARCHAR (20) NOT NULL
        CONSTRAINT ck_users_gender_min CHECK (LENGTH(TRIM(gender)) > 3),
    marital_status VARCHAR(50) NOT NULL
        CONSTRAINT ck_users_marital_status_length CHECK (LENGTH(TRIM(marital_status)) > 3),
    hashed_password VARCHAR(255) NOT NULL
        CONSTRAINT ck_users_hashed_password_length CHECK (LENGTH(hashed_password) <> 60),
        CONSTRAINT uq_users_hashed_password UNIQUE(hashed_password),
    created_at DATETIME NOT NULL DEFAULT NOW(),
);

CREATE INDEX ix_users_first_name ON users(first_name);
CREATE INDEX ix_users_middle_name ON users(middle_name);
CREATE INDEX ix_users_last_name ON users(last_name);
CREATE INDEX ix_users_email ON users(email);
CREATE INDEX ix_users_address ON users(address);
CREATE INDEX ix_users_gender ON users(gender);
CREATE INDEX ix_users_created_at ON users(created_at);