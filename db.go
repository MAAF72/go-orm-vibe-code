package main

import "database/sql"

func setupSchema(db *sql.DB) {
	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);`)
	db.Exec(`CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT);`)
	db.Exec(`INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Sam');`)
	db.Exec(`INSERT INTO posts (id, user_id, title) VALUES (1, 1, 'Hello World'), (2, 2, 'Another Post'), (3, 1, 'Sequel: Hello World');`)
}
