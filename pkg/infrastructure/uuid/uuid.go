package uuid

import (
	"github.com/google/uuid"
)

type Generator struct {
}

func (U *Generator) NewUUID() string {
	return uuid.New().String()
}

func NewGenerator() *Generator {
	return &Generator{}
}
