package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
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

	Persons = []Person{}
)

func CreateUser() error {
	//u := new(Person)
	name := "Anton"
	works := true
	_, err := Conn.Exec(context.Background(), "insert into person(name,works) values($1,$2)", name, works)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func SelectAllUsers() error {
	rows, err := Conn.Query(context.Background(), "select * from person")

	for rows.Next() {
		p := Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Persons = append(Persons, p)
	}
	/*for _, p := range Persons {
		fmt.Println(p.ID, p.Name, p.Works)
	}*/
	return err
}

func DeleteUser(id int) error {
	_, err := Conn.Exec(context.Background(), "delete from person where id=$1", id)
	return err
}
