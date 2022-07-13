package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, person *model.Person) error {
	_, err := r.pool.Exec(ctx, "insert into persons(name,works) values($1,$2)", &person.Name, &person.Works)
	if err != nil {
		err = fmt.Errorf("database error with create user: %v", err)
	}
	return err
}

func (r *Repository) SelectAll() ([]*model.Person, error) {
	var persons []*model.Person

	rows, err := r.pool.Query(context.Background(), "select id,name,works from persons")
	if err != nil {
		return nil, fmt.Errorf("database error with select all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		p := model.Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works)
		if err != nil {
			fmt.Println(err)
			continue
		}
		persons = append(persons, &p)
	}

	return persons, err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "delete from persons where id=$1", id)
	if err != nil {
		return fmt.Errorf("error with delete user %v", err)
	}
	return nil
}

func (r *Repository) Update(id int, person *model.Person) error {
	_, err := r.pool.Exec(context.Background(), "update persons set name=$1 where id=$2", person.Name, id)
	if err != nil {
		return fmt.Errorf("error with update user %v", err)
	}
	return err
}
func (r *Repository) SelectById(id int) (model.Person, error) {
	p := model.Person{}
	err := r.pool.QueryRow(context.Background(), "select id,name,works from persons where id=$1", id).Scan(&p.ID, &p.Name, &p.Works)
	if err != nil {
		return p, fmt.Errorf("database error %v", err)
	}
	return p, nil
}
