package useredit

import (
	"log"
	"os"
	"testing"

	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/ldap"
	"github.com/zk101/nixie/lib/testutil"
	"github.com/zk101/nixie/models/ldap/auth/usernew"
)

// Global vars used for testing
var (
	conf     *testutil.Config
	connLDAP *ldap.Client
)

// defaultUser returns a pre-populated Model
func defaultUser() *usernew.Model {
	user := usernew.New()

	user.DN = "uid=testuser,ou=users," + conf.LDAP.Base
	user.Base = conf.LDAP.Base
	user.Name = "Test User"
	user.Username = "testuser"
	user.Password = "password"

	return user
}

// TestMain sets up LDAP and tears it down again
func TestMain(m *testing.M) {
	var err error

	conf, err = testutil.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	cacert, err := config.LoadCAcerts(conf.Controls.CAcertPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	connLDAP = ldap.NewClient(&conf.LDAP, cacert)
	defer connLDAP.Close()

	if err := connLDAP.Test(); err != nil {
		log.Fatal(err.Error())
	}

	usernew := defaultUser()
	if err := usernew.Create(connLDAP); err != nil {
		log.Fatal(err.Error())
	}

	retCode := m.Run()

	if err := usernew.Remove(connLDAP); err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(retCode)
}

// TestEditName
func TestEditName(t *testing.T) {
	useredit := New()
	useredit.Base = conf.LDAP.Base
	useredit.Username = "testuser"

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	useredit.Name = "New Name"

	if err := useredit.Edit(connLDAP); err != nil {
		t.Errorf("TestEdit edit failed: %s\n", err.Error())
	}

	useredit.Name = ""

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	if useredit.Name != "New Name" {
		t.Error("TestEdit check name failed")
	}
}

// TestEditEmail
func TestEditEmail(t *testing.T) {
	useredit := New()
	useredit.Base = conf.LDAP.Base
	useredit.Username = "testuser"

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	useredit.Email = "asd@dummy.org"

	if err := useredit.Edit(connLDAP); err != nil {
		t.Errorf("TestEdit edit failed: %s\n", err.Error())
	}

	useredit.Email = ""

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	if useredit.Email != "asd@dummy.org" {
		t.Error("TestEdit check name failed")
	}

	useredit.Email = "dsa@dummy.org"

	if err := useredit.Edit(connLDAP); err != nil {
		t.Errorf("TestEdit edit failed: %s\n", err.Error())
	}

	useredit.Email = ""

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	if useredit.Email != "dsa@dummy.org" {
		t.Error("TestEdit check name failed")
	}

	useredit.Email = ""

	if err := useredit.Edit(connLDAP); err != nil {
		t.Errorf("TestEdit edit failed: %s\n", err.Error())
	}

	if err := useredit.Fetch(connLDAP); err != nil {
		t.Errorf("TestEdit fetch failed: %s\n", err.Error())
	}

	if useredit.Email != "" {
		t.Error("TestEdit check name failed")
	}
}

// EOF
