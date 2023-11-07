CREATE TABLE IF NOT EXISTS products {
    id SERIAL PRIMARY KEY,
    title VARCHAR(64),
    img VARCHAR,
    price FLOAT,
    quantity INT,
    seller_name VARCHAR(64)

}