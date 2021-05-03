package phonebook

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Paulo-Eduardo/phone_book/database"
)

func insertPhonebook(phoneBook Phonebook) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO phonebooks
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

func getPhonebookList() ([]Phonebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, 
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

func getPhonebook(phonebookID int) (*Phonebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	phonebookId,
	name,
	phone,
	email
	FROM phonebooks
	WHERE phonebookId = ?`, phonebookID)

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

func removePhonebook(phonebookID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM phonebooks where phonebookId = ?`, phonebookID)
	if err != nil {
		return err
	}
	return nil
}

func updatePhonebook(phonebook Phonebook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE phonebooks SET
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

func searchPhonebookForName(name string) ([]Phonebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, `SELECT
	phonebookId,
	name,
	phone,
	email
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