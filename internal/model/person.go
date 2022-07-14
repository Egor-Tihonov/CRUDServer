package model

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Works bool   `json:"works"`
	Age   int    `json:"age"`
}
