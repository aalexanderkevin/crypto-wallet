CREATE TABLE users (
	id VARCHAR (255) PRIMARY KEY,
	username VARCHAR (255) UNIQUE NOT NULL,
	email VARCHAR (255) UNIQUE NOT NULL,
	full_name VARCHAR (255) NOT NULL,
	password VARCHAR(200) NOT NULL,
	password_salt VARCHAR(200) NOT NULL,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);

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
