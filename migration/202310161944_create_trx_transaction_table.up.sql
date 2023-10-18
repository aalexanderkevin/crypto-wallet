CREATE TABLE trx_transactions (
	id VARCHAR(255) PRIMARY KEY,
	sender_address VARCHAR(255) NOT NULL,
	receiver_address VARCHAR(255) NOT NULL,
	amount INT NOT NULL,
	fee INT NOT NULL,
	block INT NULL,
	confirmation INT NOT NULL DEFAULT 0,
	status VARCHAR(10) NOT NULL,
	received_at timestamp NULL,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);
