package structs

import "time"

type RegisterCache struct {
	Id       string       `xorm:"not null pk VARCHAR(50)"`
	Code string    `xorm:"not null VARCHAR(10)"`
	Created time.Time `xorm:"TIMESTAMP created"`
}

type User struct {
	Id string `xorm:"not null pk VARCHAR(40)"`
	Username string `xorm:"not null unique VARCHAR(50)"`
	Password string `xorm:"not null VARCHAR(50)"`
	Code string
	Created time.Time `xorm:"TIMESTAMP created"`
}