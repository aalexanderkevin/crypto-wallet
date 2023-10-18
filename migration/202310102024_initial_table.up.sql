CREATE TABLE btc_transactions (
	id VARCHAR(255) PRIMARY KEY,
	sender_address TEXT [] NOT NULL,
	receiver_address TEXT [] NOT NULL,
	amount INT NOT NULL,
	fee INT NOT NULL,
	block INT NULL,
	confirmation INT NOT NULL,
	status VARCHAR(10) NOT NULL,
	received_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	completed_at timestamp NULL
);
