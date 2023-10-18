CREATE TABLE wallets (
	id VARCHAR(255) PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	seed_phrase bytea NOT NULL,
	btc_address VARCHAR(255) NOT NULL,
	eth_address VARCHAR(255) NOT NULL,
	trx_address VARCHAR(255) NOT NULL,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL 
);
