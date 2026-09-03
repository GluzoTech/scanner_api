CREATE TABLE scans (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    barcode_input TEXT NOT NULL,
    created_at DATE NOT NULL,
    captured_at TIME NOT NULL
);