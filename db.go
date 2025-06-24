package main

import "database/sql"

func setupSchema(db *sql.DB) {
	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, supervisor_id INTEGER);`)
	db.Exec(`CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT);`)
	db.Exec(`INSERT INTO users (id, name, supervisor_id) VALUES (1, 'Boss', NULL), (2, 'Bob', 1), (3, 'Sam', 1), (4, 'Alice', 1);`)
	db.Exec(`INSERT INTO posts (id, user_id, title) VALUES (1, 1, 'Hello World'), (2, 2, 'Another Post'), (3, 1, 'Sequel: Hello World');`)
}
