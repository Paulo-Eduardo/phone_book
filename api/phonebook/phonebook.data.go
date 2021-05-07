package phonebook

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"time"
)

func insert(phoneBook Phonebook, db *sql.DB, timeout int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	result, err := db.ExecContext(ctx, `INSERT INTO phonebooks
	(name,
	phone,
	email) VALUES (?, ?, ?)`,
		phoneBook.Name,
		phoneBook.Phone,
		phoneBook.Email)

	if err != nil {
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertID), nil
}

func get(phonebookID int, db *sql.DB, timeout int) (*Phonebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT phonebookId, name, phone, email FROM phonebooks WHERE phonebookId = ?", phonebookID)

	phonebook := &Phonebook{}
	err := row.Scan(
		&phonebook.PhonebookID,
		&phonebook.Name,
		&phonebook.Phone,
		&phonebook.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return phonebook, nil
}

func remove(phonebookID int, db *sql.DB, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `DELETE FROM phonebooks where phonebookId = ?`, phonebookID)
	if err != nil {
		return err
	}
	return nil
}

func update(phonebook Phonebook, db *sql.DB, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `UPDATE phonebooks SET
	name=?,
	phone=?,
	email=?
	WHERE phonebookId = ?`,
		phonebook.Name,
		phonebook.Phone,
		phonebook.Email,
		phonebook.PhonebookID)

	if err != nil {
		return err
	}
	return nil
}

func list(query url.Values, db *sql.DB, timeout int) ([]Phonebook, error) {
	if query["name"] != nil {
		return searchForName(query.Get("name"), db, timeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	results, err := db.QueryContext(ctx,
		`SELECT 
		phonebookId,
		name,
		email,
		phone
		FROM phonebooks`)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	phonebooks := make([]Phonebook, 0)

	for results.Next() {
		var phonebook Phonebook
		results.Scan(
			&phonebook.PhonebookID,
			&phonebook.Name,
			&phonebook.Email,
			&phonebook.Phone)

		phonebooks = append(phonebooks, phonebook)
	}

	return phonebooks, nil
}

func searchForName(name string, db *sql.DB, timeout int) ([]Phonebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	results, err := db.QueryContext(ctx, `SELECT
	phonebookId,
	name,
  email,
	phone
	FROM phonebooks
	WHERE name LIKE ?`, "%"+name+"%")

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	phonebooks := make([]Phonebook, 0)

	for results.Next() {
		var phonebook Phonebook
		results.Scan(
			&phonebook.PhonebookID,
			&phonebook.Name,
			&phonebook.Email,
			&phonebook.Phone)

		phonebooks = append(phonebooks, phonebook)
	}

	return phonebooks, nil
}
