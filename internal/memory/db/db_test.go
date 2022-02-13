package db_test

import (
	"database/sql"
	"errors"
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/internal/memory/db"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockStorage struct {
}

var checkEntry memory.MemoryEntry
var checkSuccess bool

func (d MockStorage) CreateDb() (r sql.Result, err error) {
	checkSuccess = true
	return r, errors.New("AAAAA")
}

func (d MockStorage) InsertRow(entry memory.MemoryEntry) (err error) {
	if checkEntry == entry {
		checkSuccess = true
	} else {
		checkSuccess = false
	}
	return err
}

func (d MockStorage) SelectRowByShortId(shortId string) (out memory.MemoryEntry, err error) {
	if checkEntry.ShortId == shortId {
		checkSuccess = true
	} else {
		checkSuccess = false
	}
	return out, err
}

func (d MockStorage) SelectRowByLongLink(longLink string) (out memory.MemoryEntry, err error) {
	if checkEntry.LongLink == longLink {
		checkSuccess = true
	} else {
		checkSuccess = false
	}
	return out, err
}

func TestDb(t *testing.T) {
	mock := MockStorage{}
	db.SetAdapter(mock)

	dbs := db.GetMemoryInstance()
	entry1 := memory.MemoryEntry{
		ShortId:   "aaabbbcccd",
		LongLink:  "https://google.com",
		CreatedAt: time.Now(),
		Author:    "test",
	}
	checkEntry = entry1
	err := dbs.AddEntry(entry1)
	require.Nil(t, err)
	require.True(t, checkSuccess)
	_, err = dbs.GetEntryByLongLink("https://google.com")
	require.Nil(t, err)
	require.True(t, checkSuccess)
	_, err = dbs.GetEntryByShortId("aaabbbcccd")
	require.Nil(t, err)
	require.True(t, checkSuccess)
	err = dbs.Clear()
	require.Nil(t, err)
}
