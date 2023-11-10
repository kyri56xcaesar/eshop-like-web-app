DROP TABLE IF EXISTS products;

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(64),
    img VARCHAR,
    price FLOAT,
    quantity INT,
    seller_name VARCHAR(64)

);