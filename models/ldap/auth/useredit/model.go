package useredit

// Model holds chat structure
type Model struct {
	DN       string
	Base     string
	Name     string
	Username string
	Email    string
	Mobile   string
	Address  string
}

// New returns a pointer to a New model struct
func New() *Model {
	return &Model{}
}

// EOF
