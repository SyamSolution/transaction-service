CREATE TABLE detail_transaction (
    detail_transaction_id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT,
    ticket_id int,
    ticket_type VARCHAR(50),
    country_name VARCHAR(100),
    city VARCHAR(100),
    quantity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
