package memory

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Memory struct {
	Storage
}

type Storage interface {
	AddEntry(entry MemoryEntry) error
	GetEntryByShortId(shortId string) (MemoryEntry, error)
	GetEntryByLongLink(longLink string) (MemoryEntry, error)
	Clear() error
}

type MemoryEntry struct {
	ShortId   string
	LongLink  string
	Author    string
	CreatedAt time.Time
}

// - Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func (m *MemoryEntry) GenerateUniqueShortLink(db Memory) {
	var shortId string
	var err error = nil
	for err == nil {
		shortId = gonanoid.MustGenerate(alphabet, 10)
		_, err = db.GetEntryByShortId(shortId)
	}
	m.ShortId = shortId
}
