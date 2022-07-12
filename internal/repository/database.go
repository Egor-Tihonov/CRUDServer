package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, person *model.Person) error {
	//u := new(Person)
	person.Name = "Egor"
	person.Works = true
	_, err := r.pool.Exec(ctx, "insert into persons(name,works) values($1,$2)", person.Name, person.Works)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r *Repository) SelectAll() []model.Person {
	var persons []model.Person

	rows, err := r.pool.Query(context.Background(), "select * from persons")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		p := model.Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works)
		if err != nil {
			fmt.Println(err)
			continue
		}
		persons = append(persons, p)
	}
	/*for _, p := range Persons {
		fmt.Println(p.ID, p.Name, p.Works)
	}*/
	return persons
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "delete from persons where id=$1", id)
	return err
}

func (r *Repository) Update(id int) error {
	name := "qwerty"
	_, err := r.pool.Exec(context.Background(), "update persons set name=$1 where id=$2", name, id)
	return err
}
func (r *Repository) SelectById(id int) []model.Person {
	var persons []model.Person

	rows, err := r.pool.Query(context.Background(), "select * from persons where id=$1", id)
	if err != nil {
		return nil
	}
	for rows.Next() {
		p := model.Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works)
		if err != nil {
			fmt.Println(err)
			continue
		}
		persons = append(persons, p)
	}
	/*for _, p := range Persons {
		fmt.Println(p.ID, p.Name, p.Works)
	}*/
	return persons
}
