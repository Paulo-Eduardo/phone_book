package phonebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Paulo-Eduardo/phone_book/cors"
)

const phonebookBasePath = "phonebooks"

func SetupRoutes(apiBasePath string) {
	handlePhonebooks := http.HandlerFunc(phonebooksHandler)
	handlePhonebook := http.HandlerFunc(phonebookHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, phonebookBasePath), cors.Middleware(handlePhonebooks))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, phonebookBasePath), cors.Middleware(handlePhonebook))
}

func phonebooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		var phonebookList []Phonebook
		var err error
		if query["name"] != nil {
			phonebookList, err = searchPhonebookForName(query.Get("name"))
		} else {
			phonebookList, err = getPhonebookList()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		phonebooksJson, err := json.Marshal(phonebookList)
		if err != nil {
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
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newPhonebook)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newPhonebook.PhonebookID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err:= insertPhonebook(newPhonebook)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonId, err := json.Marshal(id)
		if err != nil {
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	phonebook, err := getPhonebook(phonebookID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if phonebook == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		phonebookJSON, err := json.Marshal(phonebook)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(phonebookJSON)
	case http.MethodPut:
		var updatedPhonebook Phonebook
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedPhonebook)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedPhonebook.PhonebookID != phonebookID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updatePhonebook(updatedPhonebook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removePhonebook(phonebookID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}