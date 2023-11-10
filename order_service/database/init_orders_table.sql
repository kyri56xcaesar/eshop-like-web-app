DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    product_title VARCHAR(64),
    quantity INT,
    buyer_name VARCHAR(64),
    seller_name VARCHAR(64)

);