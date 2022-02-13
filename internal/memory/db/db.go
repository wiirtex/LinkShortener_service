package db

import (
	"database/sql"
	"errors"
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/pkg/helpers"

	_ "github.com/lib/pq"
)

type databaseAdapter interface {
	CreateDb() (sql.Result, error)
	InsertRow(entry memory.MemoryEntry) (err error)
	SelectRowByShortId(shortId string) (out memory.MemoryEntry, err error)
	SelectRowByLongLink(longLink string) (out memory.MemoryEntry, err error)
}

type databaseSender struct {
	source *sql.DB
}

func (d databaseSender) CreateDb() (sql.Result, error) {
	return d.source.Exec(`CREATE TABLE IF NOT EXISTS ozonlinks (
		id SERIAL PRIMARY KEY,
		longLink TEXT NOT NULL,
		shortId VARCHAR(10) NOT NULL,
		author TEXT NOT NULL,
		createdAt timestamp NOT NULL);
		CREATE UNIQUE INDEX IF NOT EXISTS longLinkIndex ON ozonlinks (longLink);
		CREATE UNIQUE INDEX IF NOT EXISTS shortIdIndex ON ozonlinks (shortId);`)
}

func (d databaseSender) InsertRow(in memory.MemoryEntry) (err error) {
	row := d.source.QueryRow(`INSERT INTO ozonlinks (longLink, shortId, author, createdAt) VALUES ($1, $2, $3, $4);`, in.LongLink, in.ShortId, in.Author, in.CreatedAt.Format("2006-01-02 15:04:05"))
	return row.Err()
}

func (d databaseSender) SelectRowByShortId(shortId string) (out memory.MemoryEntry, err error) {
	row := d.source.QueryRow(`select longLink, shortId, author, createdAt from ozonlinks where shortId = $1;`, shortId)

	if row.Err() != nil {
		return out, row.Err()
	}

	err = row.Scan(&out.LongLink, &out.ShortId, &out.Author, &out.CreatedAt)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return out, errors.New("No such entry")
	}
	return out, err
}

func (d databaseSender) SelectRowByLongLink(longLink string) (out memory.MemoryEntry, err error) {
	row := d.source.QueryRow(`select longLink, shortId, author, createdAt from ozonlinks where longLink = $1;`, longLink)

	if row.Err() != nil {
		return out, row.Err()
	}

	err = row.Scan(&out.LongLink, &out.ShortId, &out.Author, &out.CreatedAt)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return out, errors.New("No such entry")
	}
	return out, err
}

type database struct {
	adapter databaseAdapter
}

var dbSingleInstance database

func SetAdapter(adapter databaseAdapter) {
	dbSingleInstance = database{
		adapter: adapter,
	}
}

func GetMemoryInstance() database {
	if dbSingleInstance.adapter == nil {
		db, err := sql.Open("postgres", helpers.GetConfig().DbConnString)
		if err != nil {
			panic("Database open error: " + err.Error())
		}
		adapter := databaseSender{
			source: db,
		}
		err = adapter.source.Ping()
		if err != nil {
			panic("Database ping error: " + err.Error())
		}

		_, err = adapter.CreateDb()
		if err != nil {
			panic("Database creating table links error: " + err.Error())
		}
		dbSingleInstance = database{
			adapter: adapter,
		}
	}
	return dbSingleInstance
}

func (d database) AddEntry(entry memory.MemoryEntry) error {
	return d.adapter.InsertRow(entry)
}
func (d database) GetEntryByShortId(shortId string) (entry memory.MemoryEntry, err error) {
	return d.adapter.SelectRowByShortId(shortId)
}
func (d database) GetEntryByLongLink(longLink string) (entry memory.MemoryEntry, err error) {
	return d.adapter.SelectRowByLongLink(longLink)
}
func (d database) Clear() error {
	return nil
}
