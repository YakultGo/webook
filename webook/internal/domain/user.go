package domain

import "time"

type User struct {
	Id          int64
	Email       string
	Password    string
	NickName    string
	Phone       string
	Birthday    time.Time
	Description string
	CreateTime  time.Time
}
