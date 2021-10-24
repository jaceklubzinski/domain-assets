package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStore struct {
	db *sql.DB
}

func NewsqliteDB() (*sqliteStore, error) {
	db, err := sql.Open("sqlite3", "domainassets.db")
	if err != nil {
		return &sqliteStore{db}, err
	}
	return &sqliteStore{db}, nil
}

func (d *sqliteStore) CreateDomainsTable() error {
	TableCreate := `
	CREATE TABLE IF NOT EXISTS domains (
			Name            TEXT NOT NULL,
			RecordType      TEXT NOT NULL,
			DNSZone         TEXT NOT NULL,
			RecordProvider  TEXT NOT NULL,
			ResourceRecords TEXT NOT NULL,
			AddetAt 		TEXT NOT NULL,
			LastUpdate 		TEXT NOT NULL,
			Status 			TEXT NOT NULL
	);
	`

	_, err := d.db.Exec(TableCreate)
	if err != nil {
		return err
	}
	return nil
}
func (d *sqliteStore) CheckIfExist(n string) (exists bool, err error) {
	sqlStmt := "SELECT name from domains where name = ?"

	switch err = d.db.QueryRow(sqlStmt, n).Scan(&n); err {
	case sql.ErrNoRows:
		exists = false
	case nil:
		exists = true
	default:
		return exists, err
	}
	return exists, nil
}

func (d *sqliteStore) LastUpdate(n string) error {
	t := time.Now().String()
	s := "Active"
	stmt, err := d.db.Prepare("UPDATE domains set LastUpdate = ?, Status = ? where name = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t, s, n)
	if err != nil {
		return err
	}
	return nil
}

func (d *sqliteStore) AddRow(n, rt, dz, rp, rr string) error {
	s := "Active"
	t := time.Now().String()
	stmt, err := d.db.Prepare("REPLACE INTO domains values(?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(n, rt, dz, rp, rr, t, t, s)
	if err != nil {
		return err
	}
	return nil
}

func (d *sqliteStore) GetRow() error {
	var title, desc string
	r, err := d.db.Query("select * from helper;")
	if err != nil {
		return err
	}
	for r.Next() {
		err := r.Scan(&title, &desc)
		if err != nil {
			return err
		}
	}
	return nil
}
