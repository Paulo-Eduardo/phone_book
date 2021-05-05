package phonebook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Paulo-Eduardo/phone_book/cors"
	"github.com/Paulo-Eduardo/phone_book/logger"
)

const phonebookBasePath = "phonebooks"

var db *sql.DB
var timeout int

func SetupRoutes(apiBasePath string, dbCoon *sql.DB, to int) {
	db = dbCoon
	timeout = to
	handlePhonebooks := http.HandlerFunc(phonebooksHandler)
	handlePhonebook := http.HandlerFunc(phonebookHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, phonebookBasePath), logger.Middleware(cors.Middleware(handlePhonebooks)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, phonebookBasePath), logger.Middleware(cors.Middleware(handlePhonebook)))
}

func phonebooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var phonebookList []Phonebook
		var err error
		phonebookList, err = list(r.URL.Query(), db, timeout)
		if err != nil {
      log.Printf("An error accured trying to list user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		phonebooksJson, err := json.Marshal(phonebookList)
		if err != nil {
      log.Printf("An error accured trying parse the list of users: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(phonebooksJson)
	case http.MethodPost:
		// add a new entry in phonebook list
		var newPhonebook Phonebook
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
      log.Printf("An error accured trying to read the body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newPhonebook)
		if err != nil {
      log.Printf("An error accured trying to parse the body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newPhonebook.PhonebookID != 0 {
      log.Printf("User passed a body with ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := insert(newPhonebook, db, timeout)
		if err != nil {
      log.Printf("An error accured trying to insert the item in the database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonId, err := json.Marshal(id)
		if err != nil {
      log.Printf("An error accured trying to parse the result id: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonId)
		return
	case http.MethodOptions:
		return
	}
}

func phonebookHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "phonebooks/")
	phonebookID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
    log.Printf("An error accured trying to parse the query id: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	phonebook, err := get(phonebookID, db, timeout)

	if err != nil {
    log.Printf("An error accured trying to get the item from id: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if phonebook == nil {
    log.Printf("An error accured trying to parse the query id: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		phonebookJSON, err := json.Marshal(phonebook)
		if err != nil {
      log.Printf("An error accured trying to parse the phonebook: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(phonebookJSON)
	case http.MethodPut:
		var updatedPhonebook Phonebook
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
      log.Printf("An error accured trying to read the body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedPhonebook)
		if err != nil {
      log.Printf("An error accured trying to parse the body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedPhonebook.PhonebookID != phonebookID {
      log.Printf("An error accured, user trying to update but ID didn't match")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = update(updatedPhonebook, db, timeout)
		if err != nil {
      log.Printf("An error accured trying to update phonebook: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		remove(phonebookID, db, timeout)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
