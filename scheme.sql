CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    last_login TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE host (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL UNIQUE,
    host_type VARCHAR(255) NOT NULL,
    application_id VARCHAR(255),
    secret VARCHAR(255)
);

CREATE TABLE project (
    id SERIAL PRIMARY KEY,
    host INT NOT NULL,
    repository INT NOT NULL,
    owner INT NOT NULL,
    created_at TIMESTAMP  NOT NULL,

    FOREIGN KEY (host) REFERENCES host(id),
    FOREIGN KEY (owner) REFERENCES account(id)
);

CREATE INDEX idx_project_repository ON project(repository);

CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    external_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    project INT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (project) REFERENCES project(id)
);

CREATE INDEX idx_product_external_id ON product(external_id);

CREATE TABLE token (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL,
    expiry TIMESTAMP NOT NULL,
    owner INT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (owner) REFERENCES account(id)
);

CREATE TABLE token_access (
    token INT NOT NULL,
    product INT NOT NULL,

    PRIMARY KEY (token, product),

    FOREIGN KEY (token) REFERENCES token(id),
    FOREIGN KEY (product) REFERENCES product(id)
);

CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    account INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    sso_id INT NOT NULL,
    host INT NOT NULL,

    FOREIGN KEY (account) REFERENCES account(id),
    FOREIGN KEY (host) REFERENCES host(id)
);

CREATE TABLE permission (
    account INT NOT NULL,
    project INT NOT NULL,
    can_download BOOLEAN NOT NULL,
    can_upload BOOLEAN NOT NULL,
    can_delete BOOLEAN NOT NULL,

    PRIMARY KEY (account, project),

    FOREIGN KEY (account) REFERENCES account(id),
    FOREIGN KEY (project) REFERENCES project(id)
);

