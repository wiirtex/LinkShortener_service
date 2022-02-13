package inMemory_test

import (
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/internal/memory/inMemory"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInMemoryDb(t *testing.T) {
	entryNil := memory.MemoryEntry{
		LongLink:  "",
		ShortId:   "",
		Author:    "",
		CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	entry1 := memory.MemoryEntry{
		LongLink:  "https://google.com",
		ShortId:   "Hello",
		Author:    "127.0.0.1",
		CreatedAt: time.Now(),
	}
	entry2 := memory.MemoryEntry{
		LongLink:  "http://yandex.com",
		ShortId:   "Hello1",
		Author:    "127.0.0.1",
		CreatedAt: time.Now(),
	}

	{
		cache := inMemory.GetMemoryInstance()
		err := cache.AddEntry(entry1)
		require.Nil(t, err)
		entry, err := cache.GetEntryByShortId("Hello")
		require.Nil(t, err)
		require.Equal(t, entry, entry1)
		entry, err = cache.GetEntryByLongLink("https://google.com")
		require.Nil(t, err)
		require.Equal(t, entry, entry1)
		err = cache.Clear()
		require.Nil(t, err)
	}
	{
		cache := inMemory.GetMemoryInstance()
		err := cache.AddEntry(entry2)
		require.Nil(t, err)
		err = cache.AddEntry(entry1)
		require.Nil(t, err)
		entry, err := cache.GetEntryByShortId("Hello_")
		require.Equal(t, err.Error(), "No such entry")
		require.Equal(t, entry, entryNil)
		entry, err = cache.GetEntryByLongLink("https://google.com_")
		require.Equal(t, err.Error(), "No such entry")
		require.Equal(t, entry, entryNil)
		entry, err = cache.GetEntryByShortId("Hello")
		require.Nil(t, err)
		require.Equal(t, entry, entry1)
		entry, err = cache.GetEntryByLongLink("https://google.com")
		require.Nil(t, err)
		require.Equal(t, entry, entry1)
		entry, err = cache.GetEntryByShortId("Hello1")
		require.Nil(t, err)
		require.Equal(t, entry, entry2)
		entry, err = cache.GetEntryByLongLink("http://yandex.com")
		require.Nil(t, err)
		require.Equal(t, entry, entry2)
		err = cache.Clear()
		require.Nil(t, err)
	}
}
