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

type Service struct { //service new
	rps Repository
}

func NewService(NewRps Repository) *Service { //create
	return &Service{rps: NewRps}
}

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
			Name:     "Ivan",
			Works:    true,
			Age:      19,
			Password: "0",
		},
		{
			Name:     "query2",
			Works:    true,
			Age:      19,
			Password: "1",
		},
	}
	testNoValidData := []model.Person{
		{
			Name:     "Egor",
			Works:    false,
			Age:      18,
			Password: "3",
		},
		{
			Name:     "qwerty",
			Works:    true,
			Age:      -5,
			Password: "4",
		},
		{
			Name:     "qwerty1",
			Works:    false,
			Age:      250,
			Password: "250",
		},
	}
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range testValidData {
		_, err := rps.rps.Create(ctx, &p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.rps.Create(ctx, &p)
		require.Error(t, err, "create error")
	}
}
func TestSelectAll(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := model.Person{
		ID:       "12",
		Name:     "Andrey",
		Works:    true,
		Age:      20,
		Password: "12",
	}

	users, err := rps.rps.SelectAll(ctx)
	require.NoError(t, err, "select all: problems with select all users")
	require.Equal(t, 3, len(users), "select all: the values are`t equals")

	_, err = Pool.Exec(ctx, "insert into persons(id,name,works,age,password) values($1,$2,$3,$4,$5)", &p.ID, &p.Name, &p.Works, &p.Age, &p.Password)
	require.NoError(t, err, "select all: insert error")
	users, err = rps.rps.SelectAll(ctx)
	require.NotEqual(t, 5, len(users), "select all: the values are equals")

}

func TestSelectById(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.rps.SelectById(ctx, "12")
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.rps.SelectById(ctx, "20")
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}

func TestUpdate(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testValidData := []model.Person{
		{
			Name:  "Masha",
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
			Name:  "Egor",
			Works: false,
			Age:   120,
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
		err := rps.rps.Update(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", &p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.rps.Update(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", &p)
		require.Error(t, err, "update error")
	}
}
