package model

type Person struct { //person
	ID           string `bson,json:"id"`
	Name         string `bson,json:"name"`
	Works        bool   `bson,json:"works"`
	Age          int    `bson,json:"age"`
	Password     string `bson,json:"password"`
	RefreshToken string `bson,json:"refreshToken"`
}
type Authentication struct {
	Password string `json:"password"`
}
type RefreshTokens struct {
	RefreshToken string `json:"refreshToken"`
}
type Response struct {
	Message  string
	FileType string
	FileSize int64
}

type Config struct {
	CurrentDB     string `env:"CURRENT_DB" envDefault:"postgres"`
	Password      string `env:"PASSWORD"`
	PostgresDbUrl string `env:"POSTGRES_DB_URL"`
	MongoDbUrl    string `env:"MONGO_DB_URL"`
	RedisURL      string `env:"REDIS_DB_URL" envDefault:"localhost:6379"`
}
