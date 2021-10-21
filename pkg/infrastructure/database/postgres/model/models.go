package model

import (
	"time"
)

type Shorturl struct {
	Hash      string    `xorm:"'hash'"`
	LongURL   string    `xorm:"'long_url'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Clickdetails struct {
	Hash      string    `xorm:"'hash'"`
	IP        string    `xorm:"'ip'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
