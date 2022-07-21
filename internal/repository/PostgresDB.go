package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	//"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PRepository struct {
	Pool *pgxpool.Pool
}

//var no-records:=errors.New("no rows in result set")

//Create : insert new user into database
func (r *PRepository) Create(ctx context.Context, person *model.Person) (error, string) {
	newID := uuid.New().String()
	_, err := r.Pool.Exec(ctx, "insert into persons(id,name,works,age,password) values($1,$2,$3,$4,$5)", newID, &person.Name, &person.Works, &person.Age, &person.Password)
	if err != nil {
		return fmt.Errorf("database error with create user: %v", err), ""
	}
	return nil, newID
}

// SelectAll : Print all users(ID,Name,Works) from database
func (r *PRepository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person
	rows, err := r.Pool.Query(ctx, "select id,name,works,age from persons")
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
func (r *PRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Pool.Exec(ctx, "delete from persons where id=$1", id)
	if err != nil {
		return fmt.Errorf("error delete user %v", err)
	}
	return nil
}

//UpdateAuth : update user refreshToken by his ID
func (r *PRepository) UpdateAuth(ctx context.Context, id string, refreshToken string) error {
	_, err := r.Pool.Exec(ctx, "update persons set refreshToken=$1 where id=$2", refreshToken, id)
	if err != nil {
		return fmt.Errorf("error with update user %v", err)
	}
	return nil
}
func (r *PRepository) Update(ctx context.Context, id string, p *model.Person) error {
	_, err := r.Pool.Exec(ctx, "update persons set name=$1,works=$2,age=$3 where id=$4", &p.Name, &p.Works, &p.Age, id)
	if err != nil {
		return fmt.Errorf("error with update user %v", err)
	}
	return nil
}

//SelectById : select one user by his ID
func (r *PRepository) SelectById(ctx context.Context, id string) (model.Person, error) {
	p := model.Person{}
	err := r.Pool.QueryRow(ctx, "select id,name,works,age,password from persons where id=$1", id).Scan(&p.ID, &p.Name, &p.Works, &p.Age, &p.Password)
	if err != nil /*err==no-records*/ {
		return p, fmt.Errorf("database error, select by id: %v", err) /*p, fmt.errorf("user with this id doesnt exist")*/
	}
	return p, nil
}
func (r *PRepository) SelectByIdAuth(ctx context.Context, id string) (model.Person, error) {
	p := model.Person{}
	err := r.Pool.QueryRow(ctx, "select id,refreshToken from persons where id=$1", id).Scan(&p.ID, &p.RefreshToken)

	if err != nil /*err==no-records*/ {
		return p, fmt.Errorf("database error, select by id: %v", err) /*p, fmt.errorf("user with this id doesnt exist")*/
	}
	return p, nil
}
