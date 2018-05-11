package dbutils

const (
	shortenURLTable = `
CREATE TABLE IF NOT EXISTS shorten_url (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	original_url VARCHAR(256) NOT NULL,
	shorten_token VARCHAR(64)
);`
)
