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
	pool := Pool
	rps := New(pool)
	ctx := context.Background()
	p1 := model.Person{
		Name:  "Ivan",
		Works: true,
		Age:   19,
	}
	p2 := model.Person{
		Name:  "Egor",
		Works: false,
		Age:   18,
	}
	p3 := model.Person{
		Name:  "qwerty",
		Works: true,
		Age:   -5,
	}
	p4 := model.Person{
		Name:  "qwerty1",
		Works: false,
		Age:   250,
	}
	p5 := model.Person{
		Name:  "query2",
		Works: true,
		Age:   19,
	}

	err := rps.Create(ctx, &p1)
	require.NoError(t, err, "create: user with this name already exist")
	err = rps.Create(ctx, &p2)
	require.Error(t, err, "create: user with this name dont exist")
	err = rps.Create(ctx, &p3)
	require.Error(t, err, "create: error with user's age, it suitable")
	err = rps.Create(ctx, &p4)
	require.Error(t, err, "create: error with user's age, it suitable")
	err = rps.Create(ctx, &p5)
	require.NoError(t, err, "create: error with user's age, it not suitable")
}
func TestSelectAll(t *testing.T) {
	pool := Pool
	rps := New(pool)
	ctx := context.Background()
	users, err := rps.SelectAll(ctx)
	require.NoError(t, err, "select all: problems with select all users")
	require.Equal(t, 4, len(users), "select all: the values are`t equals")
	require.NotEqual(t, 10, len(users), "select all: the values are equals")

}

func TestSelectById(t *testing.T) {
	pool := Pool
	rps := New(pool)
	ctx := context.Background()
	_, err := rps.SelectById(ctx, 1)
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.SelectById(ctx, 20)
	require.Error(t, err, "select user by id: this id already exist")
}

func TestUpdate(t *testing.T) {
	pool := Pool
	rps := New(pool)
	ctx := context.Background()
	p1 := model.Person{
		Name:  "Egor",
		Works: true,
		Age:   19,
	}
	p2 := model.Person{
		Name:  "Artem",
		Works: false,
		Age:   18,
	}
	p3 := model.Person{
		Name:  "qwerty",
		Works: true,
		Age:   -5,
	}
	p4 := model.Person{
		Name:  "qwerty1",
		Works: false,
		Age:   250,
	}
	p5 := model.Person{
		Name:  "query21",
		Works: true,
		Age:   19,
	}

	err := rps.Update(ctx, 1, &p1)
	require.NoError(t, err, "update: user with this name already exist")
	err = rps.Update(ctx, 1, &p2)
	require.Error(t, err, "update: user with this name dont exist")
	err = rps.Update(ctx, 1, &p3)
	require.Error(t, err, "update: error with user's age, it is suitable")
	err = rps.Update(ctx, 1, &p4)
	require.Error(t, err, "update: error with user's age, it is suitable")
	err = rps.Update(ctx, 1, &p5)
	require.NoError(t, err, "update: error with user's age, it is not suitable")
	err = rps.Update(ctx, 1, &p1)
	require.NoError(t, err, "update: index isn't suitable")
	err = rps.Update(ctx, 20, &p1)
	require.NoError(t, err, "update: index is suitable")

}
