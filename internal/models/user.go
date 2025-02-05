package models

type UserLvl string

const (
	UserLvlDefault      UserLvl = "default"
	UserLvlBeginner     UserLvl = "beginner"
	UserLvlIntermediate UserLvl = "intermediate"
	UserLvlAdvanced     UserLvl = "advanced"
)

type UserEntity struct {
	ID  int64   `db:"user_id"`
	Lvl UserLvl `db:"lvl"`
}
