package userchat

/*
	Groups is a map using group name as the key, userid as the value
	Friends is a map using the friends userid as the key, and if the firendship is validanted (requested (false) or approved (true))
*/

// Model holds chat structure
type Model struct {
	DN      string
	Friends map[string]bool
}

// New returns a pointer to a New model struct
func New() *Model {
	return &Model{
		Friends: make(map[string]bool, 0),
	}
}

// EOF
