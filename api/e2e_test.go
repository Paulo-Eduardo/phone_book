package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type Phonebook struct {
	PhonebookID int    `json:phonebookId`
	Name        string `json:name`
	Phone       string `json:phone`
	Email       string `json:email`
}

func TestShouldInsertANewPhonebook(t *testing.T) {
	resp, err := CreateNewUser("Teste")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusCreated)
	}

	id, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if _, err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()
}

func TestListPhonebooks(t *testing.T) {
	resp, err := http.Get("http://localhost:5000/api/phonebooks")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}
}

func TestListPhonebooksByName(t *testing.T) {
	resp, err := CreateNewUser("Test54321")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if _, err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	resp, err = http.Get("http://localhost:5000/api/phonebooks?name=Test54321")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(respBody), "Test54321") {
		t.Error("Get didn't return name in the body")
	}
}

func TestGetPhonebooksById(t *testing.T) {
	resp, err := CreateNewUser("Test GET")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if _, err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	resp, err = http.Get(fmt.Sprintf("http://localhost:5000/api/phonebooks/%d", id))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(respBody), "Test GET") {
		t.Errorf("Get didn't return name in the body: %v", string(respBody))
	}

}

func TestPutPhonebooks(t *testing.T) {
	resp, err := CreateNewUser("Test PUT")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	id, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if _, err = DeleteUser(id); err != nil {
			t.Fatal(err)
		}
	}()

	pb := Phonebook{
		PhonebookID: id,
		Name:        "New Name",
		Email:       "newemail@t.com",
		Phone:       "1234-4321",
	}

	body, err := json.Marshal(pb)

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:5000/api/phonebooks/%v", string(respBody)), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}
}

func CreateNewUser(name string) (*http.Response, error) {
	reqBody, err := json.Marshal(map[string]string{
		"name":  name,
		"email": "t.t@t.com",
		"phone": "1234-1234",
	})

	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:5000/api/phonebooks",
		"application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DeleteUser(id int) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:5000/api/phonebooks/%d", id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
