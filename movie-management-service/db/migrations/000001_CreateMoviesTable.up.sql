CREATE TYPE genre AS ENUM('ACTION', 'COMEDY', 'DRAMA', 'FANTASY', 'HORROR', 'SCI-FI', 'THRILLER');

CREATE TABLE IF NOT EXISTS Movies (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "title" TEXT NOT NULL,
    "releaseDate" DATE,
    "genre" genre,
    "director" TEXT,
    "description" TEXT
);