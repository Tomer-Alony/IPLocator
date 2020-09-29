package request_limit

import (
	"github.com/segmentio/ksuid"
	"time"
)

// token factory function creates a new token
type tokenFactory func() *Token

type Token struct {
	ID string
	CreatedAt time.Time
}

func NewToken() *Token {
	return &Token{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
	}
}
