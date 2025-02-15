-- Create the farm_inputs database
CREATE DATABASE farm_inputs;

-- Connect to the farm_inputs database
\c farm_inputs;

-- Create the products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL,
    image_url VARCHAR(255),
    brand VARCHAR(255),
    unit VARCHAR(50),
    rating NUMERIC(3, 2),
    reviews INT
);

