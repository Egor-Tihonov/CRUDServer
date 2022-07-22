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
