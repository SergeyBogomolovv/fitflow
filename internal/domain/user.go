package domain

import "errors"

type UserLvl string

const (
	UserLvlDefault      UserLvl = "default"
	UserLvlBeginner     UserLvl = "beginner"
	UserLvlIntermediate UserLvl = "intermediate"
	UserLvlAdvanced     UserLvl = "advanced"
)

type User struct {
	ID  int64
	Lvl UserLvl
}

var ErrUserNotFound = errors.New("user not found")
