CREATE TABLE periode (
    periode_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    start_time DATE,
    end_time DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
