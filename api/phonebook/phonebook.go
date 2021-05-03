package phonebook

type Phonebook struct {
	PhonebookID int `json:phonebookId`
	Name string `json:name`
	Phone string `json:phone`
	Email string `json:email`
}