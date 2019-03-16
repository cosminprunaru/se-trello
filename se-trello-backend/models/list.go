package models

/*
BoardList structure defining what a list contained in a board looks like
*/
type BoardList struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}
