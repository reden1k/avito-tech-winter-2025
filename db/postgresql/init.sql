DO
$$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'avitostore') THEN
      CREATE DATABASE avitostore;
   END IF;
END
$$;

\c avitostore;

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    coins INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    employee_id INT,
    item_id INT,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT,
    receiver_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES employees(id),
    FOREIGN KEY (receiver_id) REFERENCES employees(id)
);
