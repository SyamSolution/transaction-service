CREATE TABLE transaction (
    transaction_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    order_id VARCHAR(50),
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_method VARCHAR(50),
    total_amount DECIMAL(10, 2),
    total_ticket INT,
    full_name VARCHAR(255),
    mobile_number VARCHAR(20),
    email VARCHAR(100),
    payment_status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
