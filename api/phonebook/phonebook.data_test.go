package phonebook

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestShouldInsertANewPhonebook(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	pb := Phonebook{
		Name:  "Nayara",
		Email: "nay.maggion@gmail.com",
		Phone: "47 996623579",
	}

	mock.ExpectExec(`INSERT INTO phonebooks`).WithArgs(
		pb.Name,
		pb.Phone,
		pb.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := insert(pb, db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}

func TestShouldGetAPhonebookById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "SELECT phonebookId, name, phone, email FROM phonebooks WHERE phonebookId = \\?"

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	if _, err := get(1, db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}

func TestShouldDeletePhonebook(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "DELETE FROM phonebooks where phonebookId = \\?"

	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := remove(1, db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}

func TestShouldUpdatePhonebook(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	pb := Phonebook{
		Name:  "Nayara",
		Email: "nay.maggion@gmail.com",
		Phone: "47 996623579",
	}

	query := "UPDATE phonebooks SET name=\\?, phone=\\?, email=\\? WHERE phonebookId = \\?"

	mock.ExpectExec(query).WithArgs(pb.Name, pb.Phone, pb.Email, pb.PhonebookID).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := update(pb, db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}

func TestShouldListPhonebooks(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "SELECT phonebookId, name, email, phone FROM phonebooks"
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com").
		AddRow("2", "Paulo Eduardo", "47996623579", "pauloes.dev@gmail.com")

	mock.ExpectQuery(query).WillReturnRows(rows)

	if _, err := list(nil, db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}

func TestShouldListPhonebooksByName(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "SELECT phonebookId, name, phone, email FROM phonebooks WHERE name LIKE \\?"

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com")

	mock.ExpectQuery(query).WithArgs("%Nay%").WillReturnRows(rows)

	if _, err := searchForName("Nay", db, 15); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
}
