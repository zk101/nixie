package presence

// GetDate Gets the date field
func (m *Model) GetDate() int64 {
	return m.cur.GetDate()
}

// GetUserid Gets the userid field
func (m *Model) GetUserid() string {
	return m.cur.GetUserid()
}

// GetDn Gets the dn field
func (m *Model) GetDn() string {
	return m.cur.GetDn()
}

// GetKey Gets the key field
func (m *Model) GetKey() string {
	return m.cur.GetKey()
}

// GetSign Gets the sign field
func (m *Model) GetSign() string {
	return m.cur.Sign
}

// GetCipher Gets the cipher field
func (m *Model) GetCipher() string {
	return m.cur.GetCipher()
}

// GetServerid Gets the serverid field
func (m *Model) GetServerid() string {
	return m.cur.GetServerid()
}

// GetConnectionid Gets the connectionid field
func (m *Model) GetConnectionid() uint64 {
	return m.cur.GetConnectionid()
}

// GetChatfriends Gets the chatfriends field
func (m *Model) GetChatfriends() map[string]bool {
	return m.cur.GetChatfriends()
}

// WOF
