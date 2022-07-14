package repository

import (
	"awesomeProject/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

var (
	Pool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatalf("Bad connection: %v", err)
	}
	Pool = pool
	run := m.Run()
	os.Exit(run)
}

func TestCreate(t *testing.T) {
	testValidData := []model.Person{
		{
			Name:  "Ivan",
			Works: true,
			Age:   19,
		},
		{
			Name:  "query2",
			Works: true,
			Age:   19,
		},
	}
	testNoValidData := []model.Person{
		{
			Name:  "Egor",
			Works: false,
			Age:   18,
		},
		{
			Name:  "qwerty",
			Works: true,
			Age:   -5,
		},
		{
			Name:  "qwerty1",
			Works: false,
			Age:   250,
		},
	}
	rps := New(Pool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range testValidData {
		err := rps.Create(ctx, &p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		err := rps.Create(ctx, &p)
		require.Error(t, err, "create error")
	}
}
func TestSelectAll(t *testing.T) {
	rps := New(Pool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := model.Person{
		Name:  "Andrey",
		Works: true,
		Age:   20,
	}

	users, err := rps.SelectAll(ctx)
	require.NoError(t, err, "select all: problems with select all users")
	require.Equal(t, 4, len(users), "select all: the values are`t equals")

	_, err = rps.pool.Exec(ctx, "insert into persons(name,works,age) values($1,$2,$3)", &p.Name, &p.Works, &p.Age)
	require.NoError(t, err, "select all: insert error")
	users, err = rps.SelectAll(ctx)
	require.NotEqual(t, 4, len(users), "select all: the values are equals")

}

func TestSelectById(t *testing.T) {
	rps := New(Pool)
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.SelectById(ctx, 1)
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.SelectById(ctx, 20)
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}

func TestUpdate(t *testing.T) {
	rps := New(Pool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testValidData := []model.Person{
		{
			Name:  "Egor",
			Works: true,
			Age:   19,
		},
		{
			Name:  "query21",
			Works: true,
			Age:   19,
		},
	}
	testNoValidData := []model.Person{
		{
			Name:  "Artem",
			Works: false,
			Age:   18,
		},
		{
			Name:  "qwerty",
			Works: true,
			Age:   -5,
		},
		{
			Name:  "qwerty1",
			Works: false,
			Age:   250,
		},
	}
	for _, p := range testValidData {
		err := rps.Update(ctx, 1, &p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.Update(ctx, 1, &p)
		require.Error(t, err, "update error")
	}
}
