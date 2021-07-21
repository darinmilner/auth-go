CREATE TABLE users (
    id bigserial PRIMARY KEY,
    name varchar NOT NULL,
    email varchar NOT NULL UNIQUE,  
	username varchar NOT NULL UNIQUE, 
	password varchar NOT NULL, 
	password_reset_code varchar DEFAULT 0,
	created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at timestamptz DEFAULT NULL
);
