package laundry

import (
	"database/sql"

	_ "github.com/jmoiron/sqlx"
)

type Booker struct {
	Id         int     `db:"id"         json:"-"`
	Identifier string  `db:"identifier" json:"identifier"` // Apartment number
	Name       *string `db:"name"       json:"name"`
	Email      *string `db:"email"      json:"email"`
	Phone      *string `db:"phone"      json:"phone"`
	Pin        *string `db:"pin"        json:"-"`
}

func (l *Laundry) GetBookers() ([]Booker, error) {
	sqlStmt := `SELECT * FROM booker`

	var bookers []Booker

	rows, err := l.db.Queryx(sqlStmt)
	if err != nil {
		l.Logger.Errorf("Could not get bookers")
		return bookers, err
	}

	defer rows.Close()

	for rows.Next() {
		var b Booker
		if err := rows.StructScan(&b); err != nil {
			l.Logger.Errorf("Could not scan row: %s", err)
			return bookers, err
		}

		bookers = append(bookers, b)
	}

	return bookers, nil
}

func (l *Laundry) GetBooker(id int) (*Booker, error) {
	sqlStmt := `SELECT * FROM booker WHERE id = ?`
	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return nil, nil
	}

	var b Booker
	if err = stmt.QueryRowx(id).StructScan(&b); err == sql.ErrNoRows {
		l.Logger.Warnf("Booker with ID %d not found", id)
		return nil, nil
	} else if err != nil {
		l.Logger.Errorf("Could not get row: %s", err)
		return nil, err
	}

	return &b, nil
}

func (l *Laundry) UpdateBooker(b *Booker) (*Booker, error) {
	sqlStmt := `
	UPDATE
		booker
	SET
		identifier = ?,
		name	   = ?,
		email      = ?,
		phone      = ?,
		pin        = ?
	WHERE
		id = ?
	`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Warnf("Could not prepare statement: %s", err)
		return nil, err
	}

	if _, err = stmt.Exec(b.Identifier, b.Name, b.Email, b.Phone, b.Pin, b.Id); err != nil {
		l.Logger.Warnf("Could not update booker with ID %d: %s", b.Id, err)
		return nil, err
	}

	return b, nil
}

func (l *Laundry) AddBooker(b *Booker) (*Booker, error) {
	return b, nil
}

func (b *Booker) Notify() *Booker {
	// Send an email

	return b
}
