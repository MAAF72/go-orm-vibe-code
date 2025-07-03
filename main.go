package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gookit/goutil/dump"
	_ "github.com/mattn/go-sqlite3"

	"github.com/maaf72/go-orm-vibe-code/model"
	"github.com/maaf72/go-orm-vibe-code/orm"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setupSchema(db)

	dump.Println(orm.Create(context.TODO(), db, model.User{ID: 5, Name: "test"}))
	dump.Println(orm.Update(context.TODO(), db, model.User{ID: 5, Name: "test updated"}))
	dump.Println(orm.Select[model.Post](context.TODO(), db))
	dump.Println(orm.Select[model.User](context.TODO(), db))
	dump.Println(orm.Get[model.User](context.TODO(), db, 1))
	dump.Println(orm.Get[model.User](context.TODO(), db, 2))

}
