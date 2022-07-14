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

//New :create new pool connection
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

//Create : insert new user into database
func (r *Repository) Create(ctx context.Context, person *model.Person) error {
	_, err := r.pool.Exec(ctx, "insert into persons(name,works,age) values($1,$2,$3)", &person.Name, &person.Works, &person.Age)
	if err != nil {
		return fmt.Errorf("database error with create user: %v", err)
	}
	return nil
}

// SelectAll : Print all users(ID,Name,Works) from database
func (r *Repository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person

	rows, err := r.pool.Query(ctx, "select id,name,works,age from persons")
	if err != nil {
		return nil, fmt.Errorf("database error with select all users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		p := model.Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Works, &p.Age)
		if err != nil {
			return nil, fmt.Errorf("database error %v", err)
		}
		persons = append(persons, &p)
	}

	return persons, err
}

//Delete : delete user by his ID
func (r *Repository) Delete(ctx context.Context, id int) error {
	id = 20
	_, err := r.pool.Exec(ctx, "delete from persons where id=$1", id)
	if err != nil {
		return fmt.Errorf("error with delete user %v", err)
	}
	return nil
}

//Update : update user by his ID
func (r *Repository) Update(ctx context.Context, id int, person *model.Person) error {
	_, err := r.pool.Exec(ctx, "update persons set name=$1,works=$2,age=$3 where id=$4", person.Name, person.Works, person.Age, id)
	if err != nil {
		return fmt.Errorf("error with update user %v", err)
	}
	return nil
}

//SelectById : select one user by his ID
func (r *Repository) SelectById(ctx context.Context, id int) (model.Person, error) {
	p := model.Person{}
	err := r.pool.QueryRow(ctx, "select id,name,works,age from persons where id=$1", id).Scan(&p.ID, &p.Name, &p.Works, &p.Age)
	if err != nil {
		return p, fmt.Errorf("database error %v", err)
	}
	return p, nil
}
