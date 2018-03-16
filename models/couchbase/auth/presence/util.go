package presence

// setupNew copies cur to new
func (m *Model) setupNew() {
	if m.new == nil {
		copyModel := *m.cur
		m.new = &copyModel
	}
}

// EOF
