package models

/*
Card structure defining what a card contained in a list looks like
*/
type Card struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
