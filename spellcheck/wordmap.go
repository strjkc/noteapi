package spellcheck

import (
	"time"
)

type WordMap struct {
	ValidWords   map[string]struct{}
	FileUpadedAt time.Time
}
