package models

type Admin struct {
	Login    string `db:"login"`
	Password []byte `db:"password"`
}
