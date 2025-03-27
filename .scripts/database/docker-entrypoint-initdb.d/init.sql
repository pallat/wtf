CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR (50),
    retail_price  NUMERIC(5,2)
);

INSERT INTO products(product_name, retail_price) VALUES ('LEATHER BELT',  '12.25');
INSERT INTO products(product_name, retail_price) VALUES ('WINTER JACKET',  '89.65');
INSERT INTO products(product_name, retail_price) VALUES ('COTTON SOCKS',  '2.85');
