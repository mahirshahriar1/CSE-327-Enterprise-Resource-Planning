CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255), -- can be NULL for new users
    role VARCHAR(50),
    department VARCHAR(50)
);
