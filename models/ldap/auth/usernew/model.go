package usernew

import uuid "github.com/satori/go.uuid"

// Model holds chat structure
type Model struct {
	DN       string
	UUID     string
	Base     string
	Name     string
	Username string
	Password string
	Email    string
	Mobile   string
	Address  string
}

// New returns a pointer to a New model struct
func New() *Model {
	return &Model{
		UUID: uuid.NewV4().String(),
	}
}

// EOF
