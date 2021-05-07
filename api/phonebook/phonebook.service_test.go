package phonebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/Paulo-Eduardo/phone_book/database"
	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	dbConn := database.New()
	SetupRoutes("", dbConn, 15)
	os.Exit(m.Run())
}
func TestPostPhonebooskHandler(t *testing.T) {
	_, rr, err := CreateNewUser("Teste")
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetPhonebooksHandler(t *testing.T) {
	handler := http.HandlerFunc(phonebooksHandler)

	req, err := http.NewRequest("GET", "/phonebooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetByNamePhonebooksHandler(t *testing.T) {
	_, rr, err := CreateNewUser("Test GetByName")
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	handler := http.HandlerFunc(phonebooksHandler)

	req, err := http.NewRequest("GET", "/phonebooks?name=Test GetByName", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`[{"PhonebookID":%d,"Name":"Test GetByName","Phone":"1234-1235","Email":"t.t@t.com"}]`, id)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPutPhoneboosHandler(t *testing.T) {
	pb, rr, err := CreateNewUser("Test PUT")
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	handler := http.HandlerFunc(phonebookHandler)

	pb.Name = "Teste 2"
	pb.PhonebookID = id

	body, err := json.Marshal(pb)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/phonebooks/%d", id), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetPhonebookHandler(t *testing.T) {
	_, rr, err := CreateNewUser("Test GET ID")
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	handler := http.HandlerFunc(phonebookHandler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/phonebooks/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"PhonebookID":%d,"Name":"Test GET ID","Phone":"1234-1235","Email":"t.t@t.com"}`, id)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeletePhonebookHandler(t *testing.T) {
	_, rr, err := CreateNewUser("Test DELETE")
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(phonebookHandler)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/phonebooks/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func CreateNewUser(name string) (Phonebook, *httptest.ResponseRecorder, error) {
	handler := http.HandlerFunc(phonebooksHandler)

	var pb = Phonebook{
		Name:  name,
		Email: "t.t@t.com",
		Phone: "1234-1235",
	}

	body, err := json.Marshal(pb)
	if err != nil {
		return pb, nil, err
	}

	req, err := http.NewRequest("POST", "/phonebooks", bytes.NewBuffer(body))
	if err != nil {
		return pb, nil, err
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	return pb, rr, nil
}

func DeleteUser(id int) error {
	handler := http.HandlerFunc(phonebookHandler)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/phonebooks/%d", id), nil)
	if err != nil {
		return err
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	return nil
}
