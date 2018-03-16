package presence

import (
	"testing"
)

// TestNew
func TestNew(t *testing.T) {
	model := New()

	if model.cur == nil {
		t.Error("New model cur pointer is nil")
	}

	if model.new != nil {
		t.Error("New model new pointer is not nil")
	}
}

// TestNewPBpresence
func TestNewPBpresence(t *testing.T) {
	presence := newPBpresence()

	if presence == nil {
		t.Error("New presence returned a nil pointer")
	}

	if presence.GetDate() <= 0 {
		t.Error("New presence date not set")
	}

	if presence.GetKey() == "" {
		t.Error("New presence key not set")
	}

	if presence.GetSign() == "" {
		t.Error("New presence sign not set")
	}

	if presence.GetCipher() == "" {
		t.Error("New presence cipher not set")
	}

	if presence.Chatfriends == nil {
		t.Fatal("New presence ChatFriends map is nil")
	}

	presence.Chatfriends["test"] = true

	if value, ok := presence.Chatfriends["test"]; ok != true || value != true {
		t.Error("New presence ChatFriends value assignment failed")
	}
}

// TestNewPBpresenceID
func TestNewPBpresenceID(t *testing.T) {
	presenceID := newPBpresenceID("test")

	if presenceID == nil {
		t.Error("New presenceID returned a nil pointer")
	}

	if presenceID.GetKey() != "test" {
		t.Error("New presenceID key not set to test")
	}
}

// TestSet
func TestSet(t *testing.T) {
	model := New()

	model.SetKey("testKey")

	if model.new.GetKey() != "testKey" {
		t.Error("model presence set key failed")
	}
}

// EOF
