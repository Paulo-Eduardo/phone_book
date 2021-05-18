package phonebook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	SetupRoutes("", db, 15)
	os.Exit(m.Run())
}

func TestPostPhonebookHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebooksHandler)

	var pb = Phonebook{
		Name:  "Create",
		Email: "t.t@t.com",
		Phone: "1234-1234",
	}

	var insertedId int64 = 2

	body, err := json.Marshal(pb)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`INSERT INTO phonebooks`).WithArgs(
		pb.Name,
		pb.Phone,
		pb.Email).WillReturnResult(sqlmock.NewResult(insertedId, 1))

	req, err := http.NewRequest("POST", "/phonebooks", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != strconv.Itoa(int(insertedId)) {
		t.Errorf("handler returned wrong inserted id: got %v want %v",
			insertedId, rr.Body.String())
	}
}

func TestGetPhonebooksHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebooksHandler)

	query := "SELECT phonebookId, name, email, phone FROM phonebooks"
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com").
		AddRow("2", "Paulo Eduardo", "47996623579", "pauloes.dev@gmail.com")

	mock.ExpectQuery(query).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/phonebooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != `[{"PhonebookID":1,"Name":"Nayara","Phone":"nay.maggioni@gmail.com","Email":"47996623579"},{"PhonebookID":2,"Name":"Paulo Eduardo","Phone":"pauloes.dev@gmail.com","Email":"47996623579"}]` {
		t.Errorf("handler returned wrong body: got %s, want %s",
			rr.Body.String(),
			`[{"PhonebookID":1,"Name":"Nayara","Phone":"nay.maggioni@gmail.com","Email":"47996623579"},
      {"PhonebookID":2,"Name":"Paulo Eduardo","Phone":"pauloes.dev@gmail.com","Email":"47996623579"}]`)
	}
}

// get filtering name

func TestGetPhonebookHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebookHandler)

	query := "SELECT phonebookId, name, phone, email FROM phonebooks WHERE phonebookId = \\?"

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/phonebooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != `{"PhonebookID":1,"Name":"Nayara","Phone":"47996623579","Email":"nay.maggioni@gmail.com"}` {
		t.Errorf("handler returned wrong body: got %s, want %s",
			rr.Body.String(),
			`{"PhonebookID":1,"Name":"Nayara","Phone":"nay.maggioni@gmail.com","Email":"47996623579"}`)
	}
}

func TestPutPhonebookHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebookHandler)

	pb := Phonebook{
		PhonebookID: 1,
		Name:        "Nayara",
		Email:       "nay.maggion@gmail.com",
		Phone:       "47 996623579",
	}

	query := "SELECT phonebookId, name, phone, email FROM phonebooks WHERE phonebookId = \\?"

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	query = "UPDATE phonebooks SET name=\\?, phone=\\?, email=\\? WHERE phonebookId = \\?"

	mock.ExpectExec(query).WithArgs(pb.Name, pb.Phone, pb.Email, pb.PhonebookID).WillReturnResult(sqlmock.NewResult(0, 1))

	pbJSON, err := json.Marshal(pb)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/phonebooks/1", bytes.NewReader(pbJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeletePhonebookHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebookHandler)

	query := "SELECT phonebookId, name, phone, email FROM phonebooks WHERE phonebookId = \\?"

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow("1", "Nayara", "47996623579", "nay.maggioni@gmail.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	query = "DELETE FROM phonebooks where phonebookId = \\?"

	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("DELETE", "/phonebooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
