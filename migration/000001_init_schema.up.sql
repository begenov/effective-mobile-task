CREATE TABLE gender (
    id SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS nationality (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS user (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50),
    age INTEGER,
    gender_id INTEGER REFERENCES gender(id),
    nationality_id INTEGER REFERENCES nationality(id)
);
