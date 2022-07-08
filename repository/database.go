package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Works bool   `json:"works"`
}

var (
	People      []*Person
	databaseUrl = ""
	Conn, _     = pgxpool.Connect(context.Background(), os.Getenv(databaseUrl))
)

func Create() error {
	//u := new(Person)
	name := "Anton"
	works := true
	_, err := Conn.Exec(context.Background(), "insert into person(name,works) values($1,$2)", name, works)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func SelectAll() {
	rows, _ := Conn.Query(context.Background(), "select * from person")
	persons := []Person{}

	for rows.Next() {
		p := Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works)
		if err != nil {
			fmt.Println(err)
			continue
		}
		persons = append(persons, p)
	}
	for _, p := range persons {
		fmt.Println(p.ID, p.Name, p.Works)
	}
}
