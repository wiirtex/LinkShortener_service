package memory_test

import (
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/internal/memory/inMemory"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	entry1 := memory.MemoryEntry{
		LongLink:  "https://google.com",
		ShortId:   "",
		Author:    "test",
		CreatedAt: time.Now(),
	}

	{
		var db = memory.Memory{
			Storage: inMemory.GetMemoryInstance(),
		}
		entry1.GenerateUniqueShortLink(db)

	}
}
