CREATE TABLE system_logs (
     id SERIAL PRIMARY KEY,
     cpu_usage FLOAT NOT NULL,
     memory_usage FLOAT NOT NULL,
     disk_usage FLOAT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);