CREATE TABLE IF NOT EXISTS application (
    id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL UNIQUE,
    pass_phrase varchar(255) NOT NULL,
    api_key varchar(64) NOT NULL UNIQUE,
    access_token_duration int DEFAULT 750 NOT NULL,
    refresh_token_duration int DEFAULT 604800 NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES application(id) NOT NULL,
    username varchar(255) NOT NULL UNIQUE,
    pass varchar(255) NOT NULL,
    last_login TIMESTAMP
);
