package model

import (
	"time"
)

type Shorturl struct {
	Hash      string    `xorm:"'hash'"`
	LongURL   string    `xorm:"'long_url'"`
	IsValid   bool      `xorm:"'is_valid'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Clickdetails struct {
	Hash      string    `xorm:"'hash'"`
	IP        string    `xorm:"'ip'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

//revive:disable:var-naming

type LoadBalancedUrlList []LoadBalancedUrl

type LoadBalancedUrl struct {
	Hash        string    `xorm:"'hash'"`
	OriginalURL string    `xorm:"'original_url'"`
	IsValid     bool      `xorm:"'is_valid'"`
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}

//revive:enable:var-naming

type DomainEvent struct {
	ID        string    `xorm:"'id'"`
	Payload   []byte    `xorm:"'payload'"`
	Type      string    `xorm:"'event_type'"`
	CreatedAt time.Time `xorm:"created"`
}
