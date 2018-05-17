package dbutils

const trainTable = `
CREATE TABLE IF NOT EXISTS train (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	driver_name VARCHAR(64) NOT NULL,
	operating_status BOOLEAN
);`

const stationTable = `
CREATE TABLE IF NOT EXISTS station (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(64) NOT NULL,
	opening_time VARCHAR(5) NULL,
	closing_time VARCHAR(5) NULL
);`

const scheduleTable = `
CREATE TABLE IF NOT EXISTS schedule (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	train_id INT NOT NULL,
	station_id INT NOT NULL,
	arrival_time VARCHAR(5) NOT NULL,
	FOREIGN KEY (train_id) REFERENCES train(id),
	FOREIGN KEY (station_id) REFERENCES station(id)
);`
