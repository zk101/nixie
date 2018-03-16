package presence

// SetDate sets the date field
func (m *Model) SetDate(value int64) {
	m.setupNew()
	m.new.Date = value
}

// SetUserid sets the userid field
func (m *Model) SetUserid(value string) {
	m.setupNew()
	m.new.Userid = value
}

// SetDn sets the dn field
func (m *Model) SetDn(value string) {
	m.setupNew()
	m.new.Dn = value
}

// SetKey sets the key field
func (m *Model) SetKey(value string) {
	m.setupNew()
	m.new.Key = value
}

// SetSign sets the sign field
func (m *Model) SetSign(value string) {
	m.setupNew()
	m.new.Sign = value
}

// SetCipher sets the cipher field
func (m *Model) SetCipher(value string) {
	m.setupNew()
	m.new.Cipher = value
}

// SetServerid sets the serverid field
func (m *Model) SetServerid(value string) {
	m.setupNew()
	m.new.Serverid = value
}

// SetConnectionid sets the connectionid field
func (m *Model) SetConnectionid(value uint64) {
	m.setupNew()
	m.new.Connectionid = value
}

// SetChatfriends sets the chatfriends field
func (m *Model) SetChatfriends(value map[string]bool) {
	m.setupNew()
	m.new.Chatfriends = value
}

// WOF
