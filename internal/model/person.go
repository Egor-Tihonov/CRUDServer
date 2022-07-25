package model

type Person struct { //person
	ID           string `bson,json:"id"`
	Name         string `bson,json:"name" validate:"required,min=6"`
	Works        bool   `bson,json:"works"`
	Age          int    `bson,json:"age" validate:"required,gte=0,lte:200"`
	Password     string `bson,json:"password" validate:"required, min=8"`
	RefreshToken string `bson,json:"refreshToken"`
}
type Authentication struct {
	Password string `json:"password"`
}
type RefreshTokens struct {
	RefreshToken string `json:"refreshToken"`
}
