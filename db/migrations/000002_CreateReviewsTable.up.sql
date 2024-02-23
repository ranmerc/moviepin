CREATE TABLE IF NOT EXISTS Reviews (
    "reviewID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "movieID" UUID REFERENCES Movies("movieID"),
    "rating" NUMERIC(2, 1) NOT NULL,
    "reviewText" TEXT,
    "createdAtUTC" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updatedAtUTC" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);