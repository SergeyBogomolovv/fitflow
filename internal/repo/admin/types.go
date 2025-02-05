package admin

type Admin struct {
	Login    string `db:"login"`
	Password []byte `db:"password"`
}
